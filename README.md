## 概览

Crab 是一个简洁优雅的 Go 微服务框架，提供生命周期管理、优雅关闭等开箱即用的能力。

## 核心理念

**简洁 · 透明 · 可控**

- ✅ **极简 API** - 仅 4 个核心方法，易学易用
- ✅ **生命周期管理** - Setup/Start/Close 三阶段清晰明确
- ✅ **优雅关闭** - 自动处理信号，安全关闭资源
- ✅ **错误上抛** - 框架不输出日志，业务层统一处理
- ✅ **类型透明** - 直接使用第三方库，无过度封装

## 安装

```bash
go get github.com/bang-go/crab
```

## 快速开始

### 基础示例

```go
package main

import (
    "github.com/bang-go/crab"
    "github.com/bang-go/crab/core/base/env"
    "github.com/bang-go/crab/core/base/logx"
)

func main() {
    app := crab.New()
    defer app.Close()  // 确保资源清理
    
    // Setup - 立即设置环境
    crab.Setup(
        func() error { return env.Build(env.WithAppEnv(env.DEV)) },
        func() error { logx.Build(logx.WithLevel(logx.LevelInfo)); return nil },
    )
    
    // Use - 注册资源生命周期
    crab.Use(crab.StartOnly(func() error {
        logx.Info("应用启动", "env", env.AppEnv())
        return nil
    }))
    
    // Run - 运行应用
    if err := app.Run(); err != nil {
        logx.Error("应用错误", "error", err)
    }
}
```

### 完整示例

```go
func main() {
    app := crab.New()
    defer app.Close()  // 确保资源清理
    
    // 🔧 Setup 阶段 - 准备环境（立即执行）
    crab.Setup(
        func() error { return env.Build(env.WithAppEnv(env.PROD)) },
        func() error { return viperx.Build(&viperx.Config{...}) },
        func() error { logx.Build(logx.WithLevel(logx.LevelInfo)); return nil },
    )
    
    // 📦 Use 阶段 - 注册资源生命周期
    
    // 数据库 - 需要启动和关闭
    var db *gorm.DB
    crab.Use(crab.StartClose(
        func() error {
            var err error
            db, err = gorm.Open(...)
            return err
        },
        func() error {
            sqlDB, _ := db.DB()
            return sqlDB.Close()
        },
    ))
    
    // 缓存预热 - 只需要启动
    crab.Use(crab.StartOnly(func() error {
        return cache.WarmUp()
    }))
    
    // HTTP 服务器
    var server *http.Server
    crab.Use(crab.StartClose(
        func() error {
            server = &http.Server{Addr: ":8080", Handler: router}
            return server.ListenAndServe()
        },
        func() error {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            return server.Shutdown(ctx)
        },
    ))
    
    // 🚀 Run 阶段 - 运行应用
    if err := app.Run(); err != nil {
        logx.Error("应用错误", "error", err)
    }
}
```

## 核心概念

### Worker 接口

Crab 提供简洁的 Worker 接口：

```go
type Worker interface {
    Setup(...types.FuncErr) error  // 立即设置环境
    Use(...Handler)                 // 注册资源生命周期
    Run() error                     // 运行应用
}
```

### Handler 构造器

Crab 提供三个构造器函数，覆盖所有场景：

```go
// 只需要启动（如缓存预热、数据加载）
func StartOnly(start types.FuncErr) Handler

// 需要启动和关闭（如数据库、服务器）
func StartClose(start, close types.FuncErr) Handler

// 只需要清理（极少见）
func CloseOnly(close types.FuncErr) Handler
```

### 生命周期

Crab 使用三阶段生命周期管理：

```
Setup（准备阶段） → Use（配置阶段） → Run（运行阶段） → Close（清理阶段）
   ↓                    ↓                  ↓                 ↓
立即执行              注册资源        触发 Start        触发 Close
环境、配置、日志    数据库、服务器    启动资源          清理资源
```

**执行时机**：
- **Setup** - 调用时立即执行，用于环境准备
- **Start** - Run() 时执行，用于启动资源
- **Close** - 信号触发或 defer 执行，用于清理资源

**重要**：推荐使用 `defer app.Close()` 确保资源清理，特别是在一次性 Job 场景下。

## 使用场景

### 场景 1: 环境、配置、日志初始化

```go
app := crab.New()
defer app.Close()  // 确保资源清理

crab.Setup(
    func() error { return env.Build(env.WithAppEnv(env.PROD)) },
    func() error { return viperx.Build(&viperx.Config{...}) },
    func() error { logx.Build(logx.WithLevel(logx.LevelInfo)); return nil },
)
```

### 场景 2: 数据库连接（需要关闭）

```go
var db *gorm.DB

app := crab.New()
defer app.Close()  // 确保数据库连接被关闭

crab.Use(crab.StartClose(
    func() error {
        var err error
        db, err = gorm.Open(...)
        return err
    },
    func() error {
        sqlDB, _ := db.DB()
        return sqlDB.Close()
    },
))
```

### 场景 3: 缓存预热（无需关闭）

```go
app := crab.New()
defer app.Close()  // 确保资源清理

crab.Use(crab.StartOnly(func() error {
    return cache.WarmUp()
}))
```

### 场景 4: HTTP 服务器

```go
var server *http.Server

app := crab.New()
defer app.Close()  // 确保资源清理

crab.Use(crab.StartClose(
    func() error {
        server = &http.Server{Addr: ":8080", Handler: router}
        return server.ListenAndServe()  // 阻塞，直到 Shutdown
    },
    func() error {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        return server.Shutdown(ctx)
    },
))

// Run() 阻塞在 ListenAndServe
// 按 Ctrl+C 触发优雅关闭
app.Run()
```

### 场景 5: 定时任务

```go
app := crab.New()
defer app.Close()  // 确保资源清理

crab.Use(crab.StartClose(
    func() error {
        cron.Start()  // 非阻塞
        return nil
    },
    func() error {
        return cron.Stop()
    },
))

app.Run()  // 立即返回

// 业务方主动阻塞
select {}  // 按 Ctrl+C 触发优雅关闭
```

## 示例

查看 [`examples/`](./examples/) 目录获取更多示例：
- [`examples/basic/`](./examples/basic/) - 基础使用示例
- [`examples/http_server/`](./examples/http_server/) - HTTP 服务器示例
- [`examples/cron_task/`](./examples/cron_task/) - 定时任务示例
- [`examples/job/`](./examples/job/) - 一次性 Job 示例

## API 参考

### 全局函数

```go
// 创建全局单例实例（使用 sync.Once 保证只创建一次）
func New() Worker

// 以下函数等同于 New().XXX()
func Setup(...types.FuncErr) error
func Use(...Handler)
func Run() error
func Close() error
```

### Handler 构造器

```go
// 创建只包含 Start 的 Handler
func StartOnly(start types.FuncErr) Handler

// 创建包含 Start 和 Close 的 Handler
func StartClose(start, close types.FuncErr) Handler

// 创建只包含 Close 的 Handler
func CloseOnly(close types.FuncErr) Handler
```

### Worker 方法

```go
// 立即设置环境（可接受多个函数）
Setup(...types.FuncErr) error

// 注册资源生命周期（可接受多个 Handler）
Use(...Handler)

// 运行应用（触发所有 Start，异步监听信号）
Run() error

// 清理资源（推荐使用 defer app.Close()）
Close() error
```

## 设计哲学

### 1. 简洁至上

- 只有 4 个核心方法，易学易用
- 使用构造器函数而非复杂配置
- 避免过度抽象和魔法

### 2. 透明可控

- 直接使用第三方库（redis.Client、gorm.DB）
- 业务方完全控制资源初始化逻辑
- 框架不隐藏实现细节

### 3. 自动化管理 + 防御性编程

- 框架自动处理信号监听
- 接收到 SIGTERM/SIGINT 时自动清理资源
- 推荐使用 `defer app.Close()` 作为安全保障（特别是一次性 Job）

### 4. 错误上抛

- 框架不输出日志
- 所有错误返回给业务层
- 业务层统一处理和记录

### 5. 机制 > 策略

- 框架提供生命周期管理机制
- 业务方实现具体策略
- 不过度封装常见场景

## License

MIT