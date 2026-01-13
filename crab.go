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
	ctx             context.Context
	cancel          context.CancelFunc
	hooks           []Hook
	shutdownTimeout time.Duration
	startupTimeout  time.Duration // 启动超时
	signals         []os.Signal
	logger          Logger // 日志接口
	mu              sync.Mutex
	state           state
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
	ctx, cancel := context.WithCancel(context.Background())
	app := &App{
		ctx:             ctx,
		cancel:          cancel,
		shutdownTimeout: 10 * time.Second,
		startupTimeout:  0, // 默认无超时
		signals:         []os.Signal{syscall.SIGTERM, syscall.SIGINT},
		state:           stateNew,
	}

	for _, opt := range opts {
		opt(app)
	}

	return app
}

// WithContext 设置基础 Context
func WithContext(ctx context.Context) Option {
	return func(a *App) {
		a.ctx, a.cancel = context.WithCancel(ctx)
	}
}

// WithShutdownTimeout 设置关闭超时时间
func WithShutdownTimeout(d time.Duration) Option {
	return func(a *App) {
		a.shutdownTimeout = d
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
	// 如果应用已经启动，禁止添加新的钩子，以确保启动顺序的确定性
	if a.state > stateNew {
		panic("crab: cannot add hook after app has started")
	}
	a.hooks = append(a.hooks, hooks...)
}

// IsRunning 返回应用是否处于运行状态（Ready）
func (a *App) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.state == stateRunning
}

// Run 启动应用并阻塞，直到收到信号或发生错误
func (a *App) Run() error {
	if !a.changeState(stateNew, stateStarting) {
		return errors.New("app already started")
	}

	a.log("App starting...")
	startBegin := time.Now()

	// 1. 启动流程 (带超时控制)
	if err := a.runStartWithTimeout(); err != nil {
		// 启动失败，执行回滚（停止已启动的组件）
		a.log("App start failed. Rolling back...", "error", err)
		_ = a.stop(context.Background())
		return err
	}

	a.log("App started successfully", "cost", time.Since(startBegin))
	a.changeState(stateStarting, stateRunning)

	// 2. 等待信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.signals...)

	select {
	case sig := <-c:
		a.log("Received signal", "signal", sig)
		signal.Stop(c)
	case <-a.ctx.Done():
		a.log("Context canceled")
	}

	// 3. 关闭流程
	return a.Stop(context.Background())
}

// Stop 手动停止应用
func (a *App) Stop(ctx context.Context) error {
	a.mu.Lock()
	if a.state == stateStopping || a.state == stateStopped {
		a.mu.Unlock()
		return nil
	}
	a.state = stateStopping
	a.mu.Unlock()

	a.log("App stopping...")
	a.cancel() // 取消主 Context

	// 创建带超时的 context 用于停止流程
	shutdownCtx, cancel := context.WithTimeout(ctx, a.shutdownTimeout)
	defer cancel()

	return a.stop(shutdownCtx)
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

// runStartWithTimeout 包装启动流程，支持超时
func (a *App) runStartWithTimeout() error {
	if a.startupTimeout > 0 {
		ctx, cancel := context.WithTimeout(a.ctx, a.startupTimeout)
		defer cancel()

		done := make(chan error, 1)
		go func() {
			done <- a.start(ctx) // 将带超时的 Context 传递给 Hook
		}()

		select {
		case err := <-done:
			return err
		case <-ctx.Done():
			return fmt.Errorf("app startup timed out after %v", a.startupTimeout)
		}
	}
	return a.start(a.ctx)
}

func (a *App) start(ctx context.Context) error {
	for i, hook := range a.hooks {
		name := hook.Name
		if name == "" {
			name = fmt.Sprintf("hook#%d", i)
		}

		// 检查超时
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if hook.OnStart != nil {
			a.log("Starting component...", "name", name)
			start := time.Now()
			if err := safeCall(ctx, hook.OnStart); err != nil {
				return fmt.Errorf("failed to start [%s]: %w", name, err)
			}
			a.log("Started component", "name", name, "cost", time.Since(start))
		}

		a.mu.Lock()
		a.hooks[i].OnStop = hook.OnStop // 确保 stop 逻辑存在
		a.mu.Unlock()
	}
	return nil
}

func (a *App) stop(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var errs []error
	for i := len(a.hooks) - 1; i >= 0; i-- {
		hook := a.hooks[i]
		name := hook.Name
		if name == "" {
			name = fmt.Sprintf("hook#%d", i)
		}

		if hook.OnStop != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("shutdown aborted: %w", ctx.Err())
			}

			a.log("Stopping component...", "name", name)
			start := time.Now()
			if err := safeCall(ctx, func(c context.Context) error { return hook.OnStop(c) }); err != nil {
				a.err("Failed to stop component", "name", name, "error", err)
				errs = append(errs, fmt.Errorf("[%s] stop failed: %w", name, err))
			} else {
				a.log("Stopped component", "name", name, "cost", time.Since(start))
			}
		}
	}

	a.state = stateStopped
	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	a.log("App stopped")
	return nil
}

func safeCall(ctx context.Context, fn types.Runner) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v\nstack: %s", r, debug.Stack())
		}
	}()
	return fn(ctx)
}

func (a *App) log(msg string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Info(context.Background(), msg, args...)
	}
}

func (a *App) err(msg string, args ...interface{}) {
	if a.logger != nil {
		a.logger.Error(context.Background(), msg, args...)
	}
}

// 全局默认实例
var std = New()

func Add(hooks ...Hook) {
	std.Add(hooks...)
}

func Run(hooks ...Hook) error {
	if len(hooks) > 0 {
		std.Add(hooks...)
	}
	return std.Run()
}

func IsRunning() bool {
	return std.IsRunning()
}

// SetLogger sets the logger for the default instance.
func SetLogger(l Logger) {
	std.logger = l
}

// Reset resets the default app instance.
// This is primarily intended for testing purposes to allow multiple Run calls.
func Reset() {
	std = New()
}
