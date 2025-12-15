package main

import (
	"context"
	"net/http"
	"time"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
)

func main() {
	app := crab.New()
	defer app.Close() // 确保资源清理

	// Setup 阶段 - 准备环境
	crab.Setup(
		func() error {
			return env.Build(env.WithAppEnv(env.DEV))
		},
		func() error {
			logx.Build(logx.WithLevel(logx.LevelInfo))
			return nil
		},
	)

	logx.Info("HTTP 服务器示例", "env", env.AppEnv())

	// HTTP 服务器 - 需要启动和优雅关闭
	var server *http.Server
	crab.Use(crab.StartClose(
		func() error {
			mux := http.NewServeMux()
			mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("pong"))
			})

			server = &http.Server{Addr: ":8080", Handler: mux}
			logx.Info("启动 HTTP 服务器", "addr", ":8080")

			// ListenAndServe 会阻塞，直到 Shutdown 被调用
			return server.ListenAndServe()
		},
		func() error {
			if server != nil {
				logx.Info("优雅关闭 HTTP 服务器")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				return server.Shutdown(ctx)
			}
			return nil
		},
	))

	// Run - 运行应用
	// 阻塞在 ListenAndServe，按 Ctrl+C 触发优雅关闭
	if err := app.Run(); err != nil {
		logx.Error("应用错误", "error", err)
	}
}
