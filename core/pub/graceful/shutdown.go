package graceful

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bang-go/crab/core/base/types"
	"github.com/bang-go/crab/core/pub/bag"
)

const (
	MaxWaitTime time.Duration = 60 * time.Second
)

// ShutdownCallback 优雅关闭时的回调函数，用于业务层记录日志
type ShutdownCallback func(event string, data map[string]any)

var shutdownBag = bag.NewBagger()
var shutdownCallback ShutdownCallback

// SetShutdownCallback 设置关闭回调函数，业务层可以通过这个回调记录日志
func SetShutdownCallback(cb ShutdownCallback) {
	shutdownCallback = cb
}

func Register(f ...types.FuncErr) {
	shutdownBag.Register(f...)
}

// WatchSignal 监听系统信号，触发优雅关闭
// 返回关闭过程中收集的所有错误
func WatchSignal(done chan struct{}, extBagger ...bag.Bagger) []error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	var reason string
	select {
	case <-done:
		reason = "application_done"
		if shutdownCallback != nil {
			shutdownCallback("shutdown_start", map[string]any{"reason": "application done signal received"})
		}
	case s := <-sigChan:
		reason = s.String()
		if shutdownCallback != nil {
			shutdownCallback("shutdown_start", map[string]any{"reason": "received signal", "signal": s.String()})
		}
	}

	return gracefulShutdown(sigChan, reason, append(extBagger, shutdownBag)...)
}

// gracefulShutdown 执行优雅关闭流程
// 返回关闭过程中的所有错误
func gracefulShutdown(sig chan os.Signal, reason string, bagger ...bag.Bagger) []error {
	signal.Stop(sig)
	ch := make(chan []error, 1)

	go func() {
		var errs []error
		for i, b := range bagger {
			if err := b.Finish(); err != nil {
				errs = append(errs, fmt.Errorf("bagger[%d] error: %w", i, err))
			}
		}
		ch <- errs
	}()

	var shutdownErrors []error
	select {
	case errs := <-ch:
		shutdownErrors = errs
		if shutdownCallback != nil {
			shutdownCallback("shutdown_complete", map[string]any{
				"reason":      reason,
				"error_count": len(errs),
			})
		}
	case <-time.After(MaxWaitTime):
		shutdownErrors = []error{errors.New("graceful shutdown timeout exceeded")}
		if shutdownCallback != nil {
			shutdownCallback("shutdown_timeout", map[string]any{
				"reason":        reason,
				"max_wait_time": MaxWaitTime.String(),
			})
		}
	}

	pro, _ := os.FindProcess(syscall.Getpid())
	_ = pro.Kill()

	return shutdownErrors
}
