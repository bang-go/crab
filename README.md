# Crab

Crab 是一个轻量级、企业级的 Go 应用生命周期管理框架。它提供了确定性的启动/关闭流程、依赖顺序管理、以及完整的可观测性支持，旨在作为构建微服务和云原生应用的稳健脚手架。

## ✨ 核心特性

*   **确定性生命周期**：
    *   **启动 (FIFO)**：严格按照 `Add` 顺序同步执行，确保配置 -> 数据库 -> 服务的依赖初始化顺序。
    *   **关闭 (LIFO)**：严格按照逆序关闭，确保上层服务先停止，底层资源后释放。
*   **企业级可观测性**：
    *   **结构化日志集成**：零适配器兼容 `slog` 及主流微服务框架日志接口，记录启动/停止耗时、组件名称等关键信息。
    *   **启动超时控制**：支持设置全局启动超时 (`WithStartupTimeout`)，防止应用初始化死锁或挂起。
    *   **组件耗时统计**：自动追踪并打印每个组件的启动/停止耗时，快速定位慢启动问题。
*   **健壮性与安全**：
    *   **自动回滚**：启动失败自动逆序清理已申请的资源。
    *   **Panic 隔离**：内置 Recover 机制，防止单组件崩溃导致进程退出。
    *   **状态保护**：应用启动后自动锁定 Hook 列表，防止运行时竞态。
*   **云原生友好**：
    *   **健康检测**：提供 `IsRunning()` 接口，用于 K8S Readiness Probe。
    *   **优雅停机**：监听系统信号，支持关闭超时控制。

## 📦 安装

```bash
go get github.com/bang-go/crab
```

## 🚀 快速开始

### 基础用法

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/bang-go/crab"
)

func main() {
	// 1. 初始化应用 (支持 Options)
	app := crab.New(
		crab.WithStartupTimeout(5*time.Second),  // 启动超时
		crab.WithShutdownTimeout(10*time.Second), // 关闭超时
	)

	// 2. 注册组件 (按依赖顺序)
	
	// 组件 A: 配置加载
	app.Add(crab.Hook{
		Name: "Config",
		OnStart: func(ctx context.Context) error {
			// Load config...
			return nil
		},
	})

	// 组件 B: HTTP 服务 (依赖配置)
	var server *http.Server
	app.Add(crab.Hook{
		Name: "HTTPServer",
		OnStart: func(ctx context.Context) error {
			server = &http.Server{Addr: ":8080"}
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	// 3. 运行
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### 集成日志与可观测性

Crab 的 `Logger` 接口设计兼容 `slog` 和主流框架（如 `bang-go/micro`）：

```go
type Logger interface {
    Info(ctx context.Context, msg string, args ...interface{})
    Error(ctx context.Context, msg string, args ...interface{})
}
```

**集成示例：**

```go
// 假设这是您的业务 Logger (例如封装了 slog)
logger := myLogger.New()

app := crab.New(
    crab.WithLogger(logger), // 直接注入，无需适配器
)

// 启动时控制台将输出结构化日志：
// [INFO] Starting component... name=HTTPServer
// [INFO] Started component name=HTTPServer cost=50ms
```

### K8S 健康检测集成

利用 `IsRunning()` 实现准确的 Readiness Probe：

```go
http.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
    if crab.IsRunning() {
        w.WriteHeader(200)
        w.Write([]byte("ok"))
    } else {
        w.WriteHeader(503)
    }
})
```

### 高级用法：依赖注入与生命周期管理

对于使用依赖注入框架（如 Google Wire, Uber Fx）的大型应用，Crab 推荐使用 **"生命周期注入模式" (Lifecycle Injection)**。

这种模式让 Provider 自行注册关闭逻辑，彻底消除了繁琐的 `Cleanup` 返回值和中间结构体。

**1. Provider 写法 (internal/provider):**

```go
// 只需要依赖 crab.Lifecycle 接口，不需要返回 cleanup 函数
func NewResource(lc crab.Lifecycle) (*Resource, error) {
    res := initResource(...) // 初始化资源
    
    // 一行代码，自动管理生命周期
    lc.Append(crab.Close(res.Close))
    
    return res, nil
}
```

**2. Main 写法 (cmd/main.go):**

```go
func main() {
    // 注入层会自动创建一个 crab.Registry 并传递给所有 Provider
    // initApp 返回收集满 Hooks 的 registry
    app, registry, err := initApp() 

    // 一键启动所有自动注册的组件
    if err := crab.Run(registry.Hooks()...); err != nil {
        log.Fatal(err)
    }
}
```

## ⚙️ 配置选项 (Options)

| Option | 说明 | 默认值 |
|--------|------|--------|
| `WithStartupTimeout(d)` | 应用启动最大允许耗时，超时则回滚 | 0 (无超时) |
| `WithShutdownTimeout(d)` | 优雅关闭最大等待时间 | 10s |
| `WithLogger(l)` | 注入日志接口，开启内部日志输出 | nil (静默) |
| `WithContext(ctx)` | 设置应用根 Context | context.Background() |
| `WithSignals(sigs...)` | 设置监听的系统信号 | SIGINT, SIGTERM |

## 💡 最佳实践

1.  **依赖顺序**：始终按照 `配置 -> 基础设施(DB/Redis) -> 业务服务 -> 对外接口` 的顺序注册 Hook。
2.  **闭包取值**：在 `OnStart` 内部读取配置值，而不是在 `Add` 时读取，以确保配置已加载（延迟求值）。
3.  **命名组件**：为每个 Hook 设置 `Name`，以便在日志中快速定位启动慢的组件。
4.  **业务 Context**：Crab 仅管理应用生命周期，业务请求的 Context（TraceID 等）应由 Web/RPC 框架处理。

## 🤝 贡献

欢迎提交 Issue 和 PR！

## 📄 许可证

MIT
