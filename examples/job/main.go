package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bang-go/crab"
)

func main() {
	app := crab.New()

	// Setup 阶段
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Setting up environment...")
			return nil
		},
	})

	log.Println("一次性 Job 示例")

	// 一次性任务
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("开始处理数据")

			// 模拟数据处理
			for i := 1; i <= 5; i++ {
				select {
				case <-ctx.Done(): // 支持上下文取消
					return ctx.Err()
				default:
					fmt.Printf("处理进度: %d/5\n", i)
					time.Sleep(500 * time.Millisecond)
				}
			}

			log.Println("数据处理完成")
			// 对于 Job 类型，如果执行完就需要退出，可以调用 Stop
			// 但通常 crab.Run 是阻塞直到信号。
			// 如果是纯 Job，可以在这里调用 app.Stop，或者手动发送信号。
			// 更好的方式可能是：Job 模式下，OnStart 运行完，如果没有其他长期运行的服务，
			// 主程序可能需要一种机制知道"所有任务已完成"。
			// 但作为一个通用的生命周期管理器，crab 默认是守护进程模式（Daemon）。
			// 如果要实现 Job 模式，可以在 OnStart 最后调用 app.Stop(context.Background())

			go func() {
				// 模拟任务完成后自动退出
				time.Sleep(100 * time.Millisecond)
				app.Stop(context.Background())
			}()

			return nil
		},
	})

	// Run - 执行任务
	if err := app.Run(); err != nil {
		log.Printf("任务执行失败: %v", err)
		return
	}

	log.Println("Job 完成")
}
