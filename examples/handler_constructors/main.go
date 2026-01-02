package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bang-go/crab"
)

func main() {
	// 实例化一个新的 App
	app := crab.New()

	// 环境设置
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Setting up environment...")
			return nil
		},
	})

	log.Println("应用启动演示")

	// 方式 1: 完整的 OnStart 和 OnStop
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("启动数据库连接")
			// db.Connect()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("关闭数据库连接")
			// db.Close()
			return nil
		},
	})

	// 方式 2: 只有 OnStart
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("缓存预热")
			// cache.WarmUp()
			return nil
		},
	})

	// 方式 3: 多个 Hook
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("加载初始数据")
			return nil
		},
	})
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("启动后台任务")
			go func() {
				// background task
			}()
			return nil
		},
	})

	// 方式 4: 闭包捕获变量（适合复杂场景）
	var server interface{ Shutdown() error }
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("启动 HTTP 服务器")
			// server = startHTTPServer()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if server != nil {
				log.Println("关闭 HTTP 服务器")
				return server.Shutdown()
			}
			return nil
		},
	})

	// 方式 5: 只有 OnStop (清理临时文件)
	var tempFile = "/tmp/app.tmp"
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Printf("创建临时文件: %s\n", tempFile)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Printf("删除临时文件: %s\n", tempFile)
			return nil
		},
	})

	// 运行应用
	if err := app.Run(); err != nil {
		log.Printf("应用错误: %v", err)
	}
}
