package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bang-go/crab"
)

// 模拟 OpenTelemetry 的 TracerProvider
type TracerProvider struct{}

func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	fmt.Println("[Tracer] 正在上报剩余的 Trace 数据...")
	time.Sleep(500 * time.Millisecond) // 模拟耗时
	fmt.Println("[Tracer] Shutdown 完成")
	return nil
}

func InitTracer() *TracerProvider {
	fmt.Println("[Tracer] 初始化全局 TracerProvider...")
	return &TracerProvider{}
}

func main() {
	// 1. 创建 App
	app := crab.New()

	var tp *TracerProvider

	// 2. 注册 Trace 初始化 Hook
	// 核心点：Crab 负责 Trace SDK 的生命周期（初始化和销毁）
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			tp = InitTracer()
			// 这里通常会调用 otel.SetTracerProvider(tp)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if tp != nil {
				// 核心点：在应用退出时，确保 Trace 数据被 Flush
				return tp.Shutdown(ctx)
			}
			return nil
		},
	})

	// 3. 注册 HTTP 服务 Hook (模拟 ginx/grpcx)
	// 核心点：业务请求的 Context 由框架中间件处理，与 Crab 的 Context 无直接关系
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			// 模拟启动 HTTP 服务
			go func() {
				fmt.Println("[Server] HTTP 服务启动 :8080")
				// 模拟一个请求处理
				handleRequest()
			}()
			return nil
		},
	})

	// 运行应用
	// 模拟运行一会后退出
	go func() {
		time.Sleep(2 * time.Second)
		app.Stop(context.Background())
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// 模拟业务请求处理
func handleRequest() {
	// 这是一个典型的业务请求 Context
	// 它通常由 ginx/grpcx 的中间件创建，包含 TraceID
	reqCtx := context.Background()

	// 模拟从 Context 获取 TraceID
	fmt.Printf("[Business] 处理请求，Ctx: %v, TraceID: <generated_by_middleware>\n", reqCtx)

	// 业务逻辑...
	time.Sleep(100 * time.Millisecond)

	fmt.Println("[Business] 请求处理完成")

	// 这里不需要 Crab 介入，Context 会随着请求结束而结束
}
