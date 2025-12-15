package crab

import (
	"fmt"
	"sync"

	"github.com/bang-go/crab/core/base/types"
	"github.com/bang-go/crab/core/pub/bag"
	"github.com/bang-go/crab/core/pub/graceful"
)

type Handler struct {
	Pre   types.FuncErr
	Init  types.FuncErr
	Close types.FuncErr
}

type Worker interface {
	Use(...Handler) error
	Run() error
	Close() error
	Done()
}

type crab struct {
	initBagger  bag.Bagger
	closeBagger bag.Bagger
	done        chan struct{}
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
			initBagger:  bag.NewBagger(),
			closeBagger: bag.NewBagger(),
			done:        make(chan struct{}, 1),
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
	// 异步监听关闭信号
	go graceful.WatchSignal(c.done, c.closeBagger)

	// 执行 Init 阶段（Pre 已在 Use() 时执行）
	// 如果业务的 Init 是阻塞的（如 gin.Run()），这里会阻塞
	// 如果业务的 Init 是非阻塞的（如 Job），执行完就返回
	if err := c.initBagger.Finish(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return get().Close()
}

func (c *crab) Close() error {
	if err := c.closeBagger.Finish(); err != nil {
		return fmt.Errorf("close bagger failed: %w", err)
	}
	return nil
}

func Use(handlers ...Handler) error {
	return get().Use(handlers...)
}

func (c *crab) Use(handlers ...Handler) error {
	for _, handler := range handlers {
		// Pre 阶段：立即执行（用于环境检查、配置加载等前置条件）
		if handler.Pre != nil {
			if err := handler.Pre(); err != nil {
				return fmt.Errorf("pre handler failed: %w", err)
			}
		}
		// Init 阶段：注册到 initBagger，在 Start() 时执行
		if handler.Init != nil {
			c.initBagger.Register(handler.Init)
		}
		// Close 阶段：注册到 closeBagger，在优雅关闭时执行
		if handler.Close != nil {
			c.closeBagger.Register(handler.Close)
		}
	}
	return nil
}

func Done() {
	get().Done()
}

func (c *crab) Done() {
	c.done <- struct{}{}
}
