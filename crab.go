package crab

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/bang-go/crab/pkg/types"
)

// Logger 定义日志接口 (兼容 bang-go/micro/logger)
type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
}

// Hook 定义应用生命周期中的一个钩子
type Hook struct {
	Name    string // 组件名称，用于日志标识
	OnStart types.Runner
	OnStop  types.Stopper
}

// Option 定义配置选项
type Option func(*App)

// App 应用实例
type App struct {
	id                        string // 应用唯一标识
	ctx                       context.Context
	cancel                    context.CancelFunc
	shutdownTimeout           time.Duration
	bestEffortShutdownTimeout time.Duration
	startupTimeout            time.Duration // 启动超时
	signals                   []os.Signal
	logger                    Logger // 日志接口

	mu                sync.RWMutex
	hooks             []Hook
	startedHooks      []Hook
	state             state
	shutdownCallbacks []func() // shutdown回调函数
	stopErr           error
	stopDone          chan struct{}
}

type state int

const (
	stateNew state = iota
	stateStarting
	stateRunning
	stateStopping
	stateStopped
)

// New 创建一个新的应用实例
func New(opts ...Option) *App {
	app := &App{
		id:                        generateAppID(),
		shutdownTimeout:           10 * time.Second,
		bestEffortShutdownTimeout: time.Second,
		startupTimeout:            0,
		signals:                   []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		state:                     stateNew,
		startedHooks:              make([]Hook, 0),
		shutdownCallbacks:         make([]func(), 0),
	}
	app.ctx, app.cancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(app)
	}

	if err := globalShutdown.Register(app); err != nil {
		app.err(app.ctx, "Failed to register app to global shutdown manager", "error", err)
	}

	return app
}

// WithContext 设置基础 Context
func WithContext(ctx context.Context) Option {
	return func(a *App) {
		if a.cancel != nil {
			a.cancel()
		}
		a.ctx, a.cancel = context.WithCancel(ctx)
	}
}

// WithShutdownTimeout 设置关闭超时时间
func WithShutdownTimeout(d time.Duration) Option {
	return func(a *App) {
		a.shutdownTimeout = d
	}
}

// WithBestEffortShutdownTimeout 设置 shutdown deadline 到期后单个 Hook 的补偿清理时限
func WithBestEffortShutdownTimeout(d time.Duration) Option {
	return func(a *App) {
		a.bestEffortShutdownTimeout = d
	}
}

// WithStartupTimeout 设置启动超时时间
func WithStartupTimeout(d time.Duration) Option {
	return func(a *App) {
		a.startupTimeout = d
	}
}

// WithSignals 设置监听的系统信号
func WithSignals(signals ...os.Signal) Option {
	return func(a *App) {
		a.signals = signals
	}
}

// WithLogger 设置日志接口
func WithLogger(l Logger) Option {
	return func(a *App) {
		a.logger = l
	}
}

// Add 注册一个或多个生命周期钩子
func (a *App) Add(hooks ...Hook) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.state > stateNew {
		panic("crab: cannot add hook after app has started")
	}
	a.hooks = append(a.hooks, hooks...)
}

// IsRunning 返回应用是否处于运行状态（Ready）
func (a *App) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state == stateRunning
}

// IsStopping 返回应用是否处于关闭流程中
func (a *App) IsStopping() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state == stateStopping
}

// Run 启动应用并阻塞，直到收到信号或发生错误
func (a *App) Run() error {
	if !a.changeState(stateNew, stateStarting) {
		return errors.New("app already started")
	}

	a.log(a.ctx, "App starting...")
	startBegin := time.Now()

	started, err := a.start(a.startContext())
	if err != nil {
		if a.isStoppingOrStopped() {
			stopErr := a.waitForStop(context.Background())
			if stopErr != nil {
				return errors.Join(err, stopErr)
			}
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		a.err(a.ctx, "App start failed. Rolling back...", "error", err)
		stopErr := a.finishStop(a.stopContext(context.Background()), started)
		if stopErr != nil {
			return errors.Join(err, stopErr)
		}
		return err
	}

	if !a.changeState(stateStarting, stateRunning) {
		return a.waitForStop(context.Background())
	}
	a.log(a.ctx, "App started successfully", "cost", formatCost(time.Since(startBegin)))

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.signals...)
	defer signal.Stop(c)

	select {
	case sig := <-c:
		a.log(a.ctx, "Received signal", "signal", sig)
	case <-a.ctx.Done():
		a.log(a.ctx, "Context canceled")
	}

	return a.Stop(context.Background())
}

// Stop 手动停止应用
func (a *App) Stop(ctx context.Context) error {
	started, done, alreadyStopping := a.beginStop()
	if alreadyStopping {
		if done == nil {
			return nil
		}
		select {
		case <-done:
			a.mu.RLock()
			defer a.mu.RUnlock()
			return a.stopErr
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	a.log(a.ctx, "App stopping...")
	a.cancel()
	callbacks := a.snapshotShutdownCallbacks()
	for _, callback := range callbacks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					a.err(context.Background(), "Shutdown callback panicked", "panic", r)
				}
			}()
			callback()
		}()
	}

	err := a.finishStop(a.stopContext(ctx), started)
	_ = globalShutdown.Unregister(a.id)
	return err
}

func (a *App) changeState(from, to state) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.state != from {
		return false
	}
	a.state = to
	return true
}

func (a *App) startContext() context.Context {
	if a.startupTimeout <= 0 {
		return a.ctx
	}
	ctx, _ := context.WithTimeout(a.ctx, a.startupTimeout)
	return ctx
}

func (a *App) start(ctx context.Context) ([]Hook, error) {
	hooks := a.snapshotHooks()
	started := make([]Hook, 0, len(hooks))

	for i, hook := range hooks {
		name := hookName(hook, i)
		if err := ctx.Err(); err != nil {
			return started, fmt.Errorf("startup aborted before [%s]: %w", name, err)
		}
		if hook.OnStart == nil {
			a.appendStartedHook(hook)
			started = append(started, hook)
			continue
		}

		a.log(ctx, "Starting component...", "name", name)
		start := time.Now()
		if err := safeCall(ctx, hook.OnStart); err != nil {
			return started, fmt.Errorf("failed to start [%s]: %w", name, err)
		}
		a.log(ctx, "Started component", "name", name, "cost", formatCost(time.Since(start)))
		a.appendStartedHook(hook)
		started = append(started, hook)
	}

	return started, nil
}

func (a *App) beginStop() ([]Hook, chan struct{}, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	switch a.state {
	case stateStopping:
		return nil, a.stopDone, true
	case stateStopped:
		return nil, nil, true
	}

	started := append([]Hook(nil), a.startedHooks...)
	a.state = stateStopping
	a.stopErr = nil
	a.stopDone = make(chan struct{})
	return started, a.stopDone, false
}

func (a *App) finishStop(ctx context.Context, hooks []Hook) error {
	err := a.stop(ctx, hooks)

	a.mu.Lock()
	a.stopErr = err
	a.state = stateStopped
	a.startedHooks = nil
	done := a.stopDone
	a.stopDone = nil
	a.mu.Unlock()

	if done != nil {
		close(done)
	}

	if err == nil {
		a.log(context.Background(), "App stopped")
	}
	return err
}

func (a *App) stop(ctx context.Context, hooks []Hook) error {
	var errs []error

	for i := len(hooks) - 1; i >= 0; i-- {
		hook := hooks[i]
		if hook.OnStop == nil {
			continue
		}

		name := hookName(hook, i)
		stopCtx := ctx
		cancel := func() {}
		if err := ctx.Err(); err != nil {
			stopCtx, cancel = a.bestEffortStopContext()
			err = fmt.Errorf("shutdown deadline exceeded before stopping [%s]: %w", name, err)
			a.err(stopCtx, "Shutdown deadline exceeded before component stop", "name", name, "error", err)
			errs = append(errs, err)
		}

		a.log(stopCtx, "Stopping component...", "name", name)
		start := time.Now()
		if err := runHook(stopCtx, func(c context.Context) error { return hook.OnStop(c) }); err != nil {
			cancel()
			a.err(stopCtx, "Failed to stop component", "name", name, "error", err)
			errs = append(errs, fmt.Errorf("[%s] stop failed: %w", name, err))
		} else {
			cancel()
			a.log(stopCtx, "Stopped component", "name", name, "cost", formatCost(time.Since(start)))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (a *App) stopContext(ctx context.Context) context.Context {
	if a.shutdownTimeout <= 0 {
		return ctx
	}

	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return ctx
		}
		if remaining <= a.shutdownTimeout {
			return ctx
		}
	}

	stopCtx, _ := context.WithTimeout(ctx, a.shutdownTimeout)
	return stopCtx
}

func (a *App) bestEffortStopContext() (context.Context, context.CancelFunc) {
	timeout := a.bestEffortShutdownTimeout
	if timeout <= 0 {
		timeout = a.shutdownTimeout / 4
	}
	if timeout <= 0 {
		timeout = 100 * time.Millisecond
	}
	if timeout > time.Second {
		timeout = time.Second
	}
	if timeout < 100*time.Millisecond {
		timeout = 100 * time.Millisecond
	}
	return context.WithTimeout(context.Background(), timeout)
}

func (a *App) snapshotHooks() []Hook {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return append([]Hook(nil), a.hooks...)
}

func (a *App) appendStartedHook(h Hook) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.startedHooks = append(a.startedHooks, h)
}

func (a *App) snapshotShutdownCallbacks() []func() {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return append([]func(){}, a.shutdownCallbacks...)
}

func (a *App) isStoppingOrStopped() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state == stateStopping || a.state == stateStopped
}

func (a *App) waitForStop(ctx context.Context) error {
	a.mu.RLock()
	state := a.state
	done := a.stopDone
	err := a.stopErr
	a.mu.RUnlock()

	if state == stateStopped && done == nil {
		return err
	}
	if done == nil {
		return nil
	}

	select {
	case <-done:
		return a.waitForStop(context.Background())
	case <-ctx.Done():
		return ctx.Err()
	}
}

func runHook(ctx context.Context, fn types.Runner) error {
	done := make(chan error, 1)
	go func() {
		done <- safeCall(ctx, fn)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func safeCall(ctx context.Context, fn types.Runner) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v\nstack: %s", r, debug.Stack())
		}
	}()
	return fn(ctx)
}

func hookName(h Hook, index int) string {
	if h.Name != "" {
		return h.Name
	}
	return fmt.Sprintf("hook#%d", index)
}

func (a *App) log(ctx context.Context, msg string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Info(ctx, msg, args...)
	}
}

func (a *App) err(ctx context.Context, msg string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Error(ctx, msg, args...)
	}
}

func (a *App) OnShutdown(fn func()) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shutdownCallbacks = append(a.shutdownCallbacks, fn)
}

func (a *App) GetID() string {
	return a.id
}

func (a *App) Unregister() error {
	return globalShutdown.Unregister(a.id)
}
