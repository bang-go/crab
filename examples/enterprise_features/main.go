package main

import (
	"context"
	"log"
	"time"

	"github.com/bang-go/crab"
)

// SimpleLogger 实现 crab.Logger 接口 (适配新接口签名)
type SimpleLogger struct{}

func (l *SimpleLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	log.Printf("[INFO] %s %v", msg, args)
}

func (l *SimpleLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, args)
}

func main() {
	// 配置企业级特性：启动超时 + 日志集成
	app := crab.New(
		crab.WithStartupTimeout(2*time.Second), // 设置2秒启动超时
		crab.WithLogger(&SimpleLogger{}),       // 注入日志
	)

	// 1. 正常的组件
	app.Add(crab.Hook{
		Name: "ConfigLoader",
		OnStart: func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond) // 模拟加载耗时
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	// 2. 耗时较长的组件 (演示正常启动日志)
	app.Add(crab.Hook{
		Name: "Database",
		OnStart: func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	// 3. 模拟超时的组件 (取消注释下面的代码来测试超时回滚)
	/*
		app.Add(crab.Hook{
			Name: "SlowService",
			OnStart: func(ctx context.Context) error {
				fmt.Println("SlowService starting (will timeout)...")
				select {
				case <-ctx.Done(): // 必须响应 ctx.Done() 才能感知超时
					return ctx.Err()
				case <-time.After(3 * time.Second): // 故意超过 2s
					return nil
				}
			},
		})
	*/

	// 演示自动退出
	go func() {
		time.Sleep(3 * time.Second)
		app.Stop(context.Background())
	}()

	if err := app.Run(); err != nil {
		log.Printf("Application exit with error: %v", err)
	}
}
