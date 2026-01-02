package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bang-go/crab"
)

func main() {
	// 实例化 App，配置 5 秒的优雅停机超时
	app := crab.New(crab.WithShutdownTimeout(5 * time.Second))

	// Setup 阶段 - 准备环境
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Setting up environment...")
			return nil
		},
	})

	log.Println("HTTP 服务器示例")

	// HTTP 服务器 - 需要启动和优雅关闭
	var server *http.Server
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			mux := http.NewServeMux()
			mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("pong"))
			})

			server = &http.Server{Addr: ":8080", Handler: mux}
			log.Println("启动 HTTP 服务器", "addr", ":8080")

			// 使用 goroutine 启动服务器，因为 OnStart 应该是阻塞的（按当前设计）
			// 注意：如果 crab 的设计是 OnStart 完成后才继续下一个，那么阻塞式服务需要在 goroutine 中运行
			// 或者，crab 可以支持阻塞式 OnStart（不返回直到 Shutdown）——但通常 OnStart 旨在初始化并启动
			// 对于 HTTP 服务，ListenAndServe 是阻塞的。
			// 如果我们希望 OnStart 阻塞直到服务就绪或失败，我们可以直接调用。
			// 但 HTTP 服务通常是长期运行的。
			// 如果 OnStart 阻塞，那么后续的 Hook 就无法执行。
			// 所以对于长期运行的服务，应该在 goroutine 中启动，并通过 channel 报告启动错误（如果需要）
			// 或者让 OnStart 只是 "启动"（非阻塞），真正的运行逻辑交给 goroutine。

			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("HTTP 服务器错误: %v", err)
					// 这里可以考虑调用 app.Stop() 或者传递错误到主流程
				}
			}()

			// 简单的等待一下确保 goroutine 启动（实际生产中可用 channel 等待 Listen 成功）
			time.Sleep(100 * time.Millisecond)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if server != nil {
				log.Println("优雅关闭 HTTP 服务器")
				// 使用传入的 ctx（已经包含了 ShutdownTimeout）
				return server.Shutdown(ctx)
			}
			return nil
		},
	})

	// Run - 运行应用
	if err := app.Run(); err != nil {
		log.Printf("应用错误: %v", err)
	}
}
