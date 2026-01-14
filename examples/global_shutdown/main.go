package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bang-go/crab"
)

func main() {
	fmt.Println("=== 全局shutdown管理器演示 ===")

	// 创建多个应用实例（都会自动注册到全局shutdown管理器）
	app1 := crab.New()
	app2 := crab.New()
	app3 := crab.New()

	// 为每个应用添加不同的生命周期钩子
	app1.Add(crab.Hook{
		Name: "Web服务",
		OnStart: func(ctx context.Context) error {
			fmt.Printf("[%s] Web服务启动中...\n", app1.GetID())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Printf("[%s] Web服务关闭中...\n", app1.GetID())
			time.Sleep(1 * time.Second) // 模拟关闭耗时
			return nil
		},
	})

	app2.Add(crab.Hook{
		Name: "数据库连接",
		OnStart: func(ctx context.Context) error {
			fmt.Printf("[%s] 数据库连接建立中...\n", app2.GetID())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Printf("[%s] 数据库连接关闭中...\n", app2.GetID())
			time.Sleep(2 * time.Second) // 模拟关闭耗时
			return nil
		},
	})

	app3.Add(crab.Hook{
		Name: "消息队列",
		OnStart: func(ctx context.Context) error {
			fmt.Printf("[%s] 消息队列启动中...\n", app3.GetID())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Printf("[%s] 消息队列关闭中...\n", app3.GetID())
			time.Sleep(500 * time.Millisecond) // 模拟关闭耗时
			return nil
		},
	})

	// 注册shutdown回调
	app1.OnShutdown(func() {
		fmt.Printf("[%s] 收到shutdown信号！\n", app1.GetID())
	})

	app2.OnShutdown(func() {
		fmt.Printf("[%s] 收到shutdown信号！\n", app2.GetID())
	})

	app3.OnShutdown(func() {
		fmt.Printf("[%s] 收到shutdown信号！\n", app3.GetID())
	})

	fmt.Printf("已注册的应用: %v\n", crab.GetApps())

	// 在后台启动所有应用
	go func() {
		if err := app1.Run(); err != nil {
			log.Printf("应用1错误: %v", err)
		}
	}()

	go func() {
		if err := app2.Run(); err != nil {
			log.Printf("应用2错误: %v", err)
		}
	}()

	go func() {
		if err := app3.Run(); err != nil {
			log.Printf("应用3错误: %v", err)
		}
	}()

	// 等待所有应用启动
	time.Sleep(2 * time.Second)

	fmt.Println("\n=== 触发全局shutdown ===")

	// 触发全局shutdown（所有应用会并行关闭）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := crab.Shutdown(ctx); err != nil {
		log.Printf("全局shutdown错误: %v", err)
	} else {
		fmt.Println("全局shutdown完成！")
	}

	fmt.Printf("剩余应用: %v\n", crab.GetApps())
}
