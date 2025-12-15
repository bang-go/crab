package crab

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bang-go/crab/core/base/types"
	"github.com/bang-go/crab/core/pub/bag"
)

type Handler struct {
	Start types.FuncErr
	Close types.FuncErr
}

// StartOnly 创建只包含 Start 的 Handler
func StartOnly(start types.FuncErr) Handler {
	return Handler{Start: start}
}

// StartClose 创建包含 Start 和 Close 的 Handler
func StartClose(start, close types.FuncErr) Handler {
	return Handler{Start: start, Close: close}
}

// CloseOnly 创建只包含 Close 的 Handler
func CloseOnly(close types.FuncErr) Handler {
	return Handler{Close: close}
}

type Worker interface {
	Use(...Handler)
	Setup(...types.FuncErr) error
	Run() error
	Close() error
}

type crab struct {
	startBagger bag.Bagger
	closeBagger bag.Bagger
	closeOnce   sync.Once // 确保 Close 只执行一次
}

var (
	instance *crab
	_        Worker = instance
	once     sync.Once
)

// New 创建全局单例 Crab 实例（使用 sync.Once 保证只创建一次）
func New() Worker {
	once.Do(func() {
		instance = &crab{
			startBagger: bag.NewBagger(),
			closeBagger: bag.NewBagger(),
		}
	})
	return instance
}

func get() Worker {
	return New()
}

func Run() error {
	return get().Run()
}

func (c *crab) Run() error {
	// 异步监听关闭信号（优雅关闭）
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		<-sigChan

		// 接收到信号，执行清理
		_ = c.close()

		// 退出进程
		os.Exit(0)
	}()

	// 执行 Start 阶段（Setup 已在调用时立即执行）
	// 如果业务的 Start 是阻塞的（如 HTTP 服务），这里会阻塞
	// 如果业务的 Start 是非阻塞的（如定时任务），执行完就返回
	if err := c.startBagger.Finish(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return get().Close()
}

// Close 公开方法，供业务方调用（推荐使用 defer app.Close()）
func (c *crab) Close() error {
	return c.close()
}

// close 私有方法，由信号处理器调用
func (c *crab) close() error {
	var err error
	c.closeOnce.Do(func() {
		if e := c.closeBagger.Finish(); e != nil {
			err = fmt.Errorf("close bagger failed: %w", e)
		}
	})
	return err
}

func Use(handlers ...Handler) {
	get().Use(handlers...)
}

func (c *crab) Use(handlers ...Handler) {
	for _, handler := range handlers {
		// Start 阶段：注册到 startBagger，在 Run() 时执行
		if handler.Start != nil {
			c.startBagger.Register(handler.Start)
		}
		// Close 阶段：注册到 closeBagger，在优雅关闭时执行
		if handler.Close != nil {
			c.closeBagger.Register(handler.Close)
		}
	}
}

func Setup(fns ...types.FuncErr) error {
	return get().Setup(fns...)
}

func (c *crab) Setup(fns ...types.FuncErr) error {
	for _, fn := range fns {
		if fn != nil {
			if err := fn(); err != nil {
				return fmt.Errorf("setup failed: %w", err)
			}
		}
	}
	return nil
}
