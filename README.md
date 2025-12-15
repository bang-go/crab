## 概览

crab 微服务框架，结合开源成熟类库，建立一体式应用框架

## 设计理念

**框架 = 辅助者，不是限制者**

- ✅ **暴露第三方库** - 直接使用 redis.Client、gorm.DB 等，业务方可查阅官方文档
- ✅ **提供便利** - 默认配置、连接验证、优雅关闭等开箱即用
- ✅ **不限制能力** - 不过度封装，保留第三方库的所有高级功能
- ✅ **降低成本** - 减少学习曲线，业务方专注业务逻辑

## 核心特性

- **生命周期管理** - Pre/Init/Close 三阶段 Handler 模式
- **优雅关闭** - 自动处理信号，安全关闭资源
- **错误上抛** - 框架不输出日志，错误传递到业务层统一处理
- **类型透明** - 直接暴露第三方库类型（Redis、MySQL、GORM 等）

## 安装

```bash
go get github.com/bang-go/crab
```

## 快速开始

### 基础用法

```go
package main

import (
    "github.com/bang-go/crab"
    "github.com/bang-go/crab/core/base/logx"
)

func main() {
    app := crab.New()
    
    // 使用辅助函数初始化日志
    crab.UseLogx(logx.WithLevel(logx.LevelInfo))
    
    // 启动应用
    if err := app.Run(); err != nil {
        logx.Error("应用错误", "error", err)
    }
}
```

### 使用辅助函数

```go
package main

import (
    "github.com/bang-go/crab"
    "github.com/bang-go/crab/core/base/env"
    "github.com/bang-go/crab/core/base/logx"
    "github.com/bang-go/crab/core/base/viperx"
    "github.com/bang-go/crab/core/db/redisx"
    "github.com/gin-gonic/gin"
)

func main() {
    app := crab.New()
    
    // 1. 初始化环境变量
    crab.UseEnv(env.WithAppEnv(env.PROD))
    
    // 2. 初始化配置
    crab.UseViper(&viperx.Config{
        ConfigFormat: viperx.FileFormatYaml,
        ConfigPaths:  []string{"./config"},
        ConfigNames:  []string{"app"},
    })
    
    // 3. 初始化日志
    crab.UseLogx(
        logx.WithLevel(logx.LevelInfo),
        logx.WithEncodeJson(),
    )
    
    // 4. 初始化 Redis
    var redisClient redisx.Client
    crab.UseRedis(&redisx.Options{
        Addr: "localhost:6379",
    }, &redisClient)
    
    // 5. 启动 Gin 服务
    crab.UseGin(":8080", func(engine *gin.Engine) {
        engine.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "pong"})
        })
    })
    
    // 启动应用
    if err := app.Run(); err != nil {
        logx.Error("应用错误", "error", err)
    }
}
```

## 核心概念

### Handler 三阶段

Crab 使用三阶段 Handler 模式管理应用生命周期：

```go
type Handler struct {
    Pre   func() error  // 立即执行（配置加载、环境检查）
    Init  func() error  // Run() 时执行（数据库连接、服务启动）
    Close func() error  // 优雅关闭时执行（资源清理）
}
```

**执行顺序**：
1. **Pre** - `Use()` 调用时立即执行，用于前置条件
2. **Init** - `Run()` 调用时执行，用于资源初始化
3. **Close** - 接收信号或调用 `Close()` 时执行，用于资源清理

### 辅助函数

Crab 提供开箱即用的辅助函数：

#### 基础组件

```go
// 环境变量
crab.UseEnv(env.WithAppEnv(env.PROD))  // PROD/DEV/TEST/PRE/GRAY

// 日志
crab.UseLogx(logx.WithLevel(logx.LevelInfo))

// 配置
crab.UseViper(&viperx.Config{...})
```

#### 数据库

```go
// Redis
var redisClient redisx.Client
crab.UseRedis(&redisx.Options{Addr: "localhost:6379"}, &redisClient)

// MySQL
var mysqlClient mysqlx.Client
crab.UseMySQL(&mysqlx.ClientConfig{...}, &mysqlClient)
```

#### HTTP 服务

```go
// Gin
crab.UseGin(":8080", func(engine *gin.Engine) {
    engine.GET("/", handler)
})

// 标准 HTTP
crab.UseHTTP(&http.Server{Addr: ":8080", Handler: handler})
```

### 单例模式

`New()` 使用 `sync.Once` 保证全局单例，多次调用返回同一实例：

```go
app1 := crab.New()
app2 := crab.New()
// app1 == app2 ✅
```

也可以直接使用全局函数：

```go
crab.Use(...)  // 等同于 crab.New().Use(...)
crab.Run()     // 等同于 crab.New().Run()
crab.Close()   // 等同于 crab.New().Close()
```

## 核心模块

### 数据库（Redis/MySQL）

直接暴露第三方库类型，查阅官方文档即可：

```go
// Redis - 使用 redis.Options
import "github.com/bang-go/crab/core/db/redisx"

client, _ := redisx.New(&redisx.Options{
    Addr: "localhost:6379",
    // ... 所有 redis.Options 的字段
})

// MySQL - 使用 driver.Config 或 DSN 字符串
import "github.com/bang-go/crab/core/db/mysqlx"

client, _ := mysqlx.New(&mysqlx.ClientConfig{
    DSN: &mysqlx.DSNConfig{
        User:   "root",
        Passwd: "password",
        Addr:   "localhost:3306",
        DBName: "mydb",
    },
})
```

### 日志（slog）

统一使用 Go 标准库 slog，支持自定义 caller skip：

```go
import "github.com/bang-go/crab/core/base/logx"

logx.Build(
    logx.WithLevel(logx.LevelInfo),
    logx.WithEncodeJson(),
    logx.WithSource(true),
    logx.WithCallerSkip(3),  // 自定义 caller
)

logx.Info("message", "key", "value")
```

### 配置（Viper）

```go
import "github.com/bang-go/crab/core/base/viperx"

viperx.Build(&viperx.Config{
    ConfigFormat: viperx.FileFormatYaml,
    ConfigPaths:  []string{"./config"},
    ConfigNames:  []string{"app"},
})
```

## 示例

查看 [`examples/`](./examples/) 目录获取更多示例

## License

MIT