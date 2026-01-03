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
	// 模拟 ginx/grpcx 业务框架
	// 在业务框架中，我们通常会注册一个 HTTP Server 来暴露业务接口和健康检测接口

	crab.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			// 创建一个 ServeMux (模拟业务路由)
			mux := http.NewServeMux()

			// 1. 业务接口
			mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello World"))
			})

			// 2. K8S Liveness Probe (存活检测)
			// 只要进程在，通常就返回 200。或者检查死锁等致命错误。
			mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			})

			// 3. K8S Readiness Probe (就绪检测)
			// 关键点：这里调用 crab.IsRunning() 来判断应用是否完全启动
			// 只有当 crab.Run() 中的所有 OnStart 钩子都执行完毕，状态才会变为 Running
			mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
				if crab.IsRunning() {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("ready"))
				} else {
					// 还在启动中，或者正在关闭中
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("not ready"))
				}
			})

			server := &http.Server{Addr: ":8080", Handler: mux}

			// 启动服务
			go func() {
				fmt.Println("[Server] Listening on :8080")
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("Server Error: %v", err)
				}
			}()

			// 模拟一个耗时的启动过程，方便观察 /readyz 的状态变化
			fmt.Println("[Init] 正在预热缓存 (3秒)...")
			time.Sleep(3 * time.Second)
			fmt.Println("[Init] 预热完成")

			return nil
		},
	})

	// 运行应用
	if err := crab.Run(); err != nil {
		log.Fatal(err)
	}
}
