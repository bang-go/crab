package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bang-go/crab"
)

func main() {
	app := crab.New()

	// Setup 阶段 - 立即设置环境（准备阶段）
	// 在新版 crab 中，Setup 和 Use 统一为 Add(Hook)
	// 如果只是初始化，只提供 OnStart 即可
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("1. 设置环境变量")
			return nil
		},
	})
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("2. 加载配置文件")
			return nil
		},
	})
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("3. 初始化日志")
			return nil
		},
	})

	log.Println("环境准备完成")

	// Use 阶段 - 注册资源生命周期
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("启动定时任务")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("停止定时任务")
			return nil
		},
	})

	// Run 阶段 - 运行应用（触发所有 Start）
	log.Println("应用启动")
	if err := app.Run(); err != nil {
		log.Printf("应用错误: %v", err)
	}
}
