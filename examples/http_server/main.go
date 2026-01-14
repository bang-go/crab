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

	app.OnShutdown(func() {
		log.Println("HTTP服务器收到shutdown信号")
	})

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
		Name: "HTTP服务器",
		OnStart: func(ctx context.Context) error {
			mux := http.NewServeMux()
			mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("pong"))
			})

			server = &http.Server{Addr: ":8080", Handler: mux}
			log.Println("启动 HTTP 服务器", "addr", ":8080")

			// 使用 goroutine 启动服务器，因为 ListenAndServe 是阻塞的
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("HTTP 服务器错误: %v", err)
				}
			}()

			// 简单的等待一下确保 goroutine 启动
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
