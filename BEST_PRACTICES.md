# Crab 框架最佳实践

## 1. 初始化顺序建议

推荐按照依赖关系的顺序初始化组件：

```go
func main() {
    app := crab.New()
    
    // 1️⃣ 环境变量（最先）- 其他组件可能依赖
    crab.UseEnv(env.WithAppEnv(env.PROD))
    
    // 2️⃣ 配置文件 - 可能依赖环境变量
    crab.UseViper(&viperx.Config{
        ConfigFormat: viperx.FileFormatYaml,
        ConfigPaths:  []string{"./config"},
        ConfigNames:  []string{fmt.Sprintf("app.%s", env.AppEnv())},
    })
    
    // 3️⃣ 日志 - 可能依赖配置
    crab.UseLogx(
        logx.WithLevel(getLogLevel()),
        logx.WithEncodeJson(),
    )
    
    // 4️⃣ 数据库等资源 - 依赖配置和日志
    var redisClient redisx.Client
    crab.UseRedis(&redisx.Options{
        Addr: viper.GetString("redis.addr"),
    }, &redisClient)
    
    // 5️⃣ HTTP 服务（最后）- 依赖所有资源
    crab.UseGin(":8080", func(engine *gin.Engine) {
        setupRoutes(engine, &redisClient)
    })
    
    // 启动
    if err := app.Run(); err != nil {
        logx.Error("应用启动失败", "error", err)
    }
}
```

## 2. 错误处理规范

### ✅ 正确做法：业务层统一处理错误

```go
func main() {
    app := crab.New()
    
    // Pre 阶段错误会立即返回
    if err := crab.UseViper(&viperx.Config{...}); err != nil {
        logx.Error("配置加载失败", "error", err)
        os.Exit(1)
    }
    
    // Init 阶段错误在 Run() 时返回
    crab.UseRedis(...)
    
    if err := app.Run(); err != nil {
        logx.Error("应用启动失败", "error", err)
        os.Exit(1)
    }
}
```

### ❌ 错误做法：框架内部打印日志

```go
// ❌ 不要在辅助函数内部打印日志
func UseRedis(...) error {
    client, err := redisx.New(opts)
    if err != nil {
        logx.Error("Redis 连接失败", "error", err)  // ❌ 错误
        return err
    }
    return nil
}

// ✅ 应该直接返回错误，由业务层处理
func UseRedis(...) error {
    client, err := redisx.New(opts)
    if err != nil {
        return err  // ✅ 正确
    }
    return nil
}
```

## 3. Handler 使用场景

### Pre 阶段 - 前置条件

**适用场景**：必须立即完成的初始化，失败则无法继续

```go
crab.Use(Handler{
    Pre: func() error {
        // ✅ 环境变量检查
        // ✅ 配置文件加载
        // ✅ 必需的环境准备
        return env.Build(...)
    },
})
```

### Init 阶段 - 资源初始化

**适用场景**：延迟到 Run() 时执行的初始化

```go
crab.Use(Handler{
    Init: func() error {
        // ✅ 数据库连接
        // ✅ 缓存预热
        // ✅ HTTP 服务启动（阻塞）
        return server.ListenAndServe()
    },
})
```

### Close 阶段 - 资源清理

**适用场景**：优雅关闭时的清理工作

```go
crab.Use(Handler{
    Close: func() error {
        // ✅ 关闭数据库连接
        // ✅ 停止 HTTP 服务
        // ✅ 刷新日志缓冲区
        return db.Close()
    },
})
```

## 4. 自定义 Handler

除了使用辅助函数，也可以自定义 Handler：

```go
// 示例：初始化消息队列
var mqClient *mq.Client

crab.Use(Handler{
    Init: func() error {
        client, err := mq.Connect(mq.Options{
            Addr: viper.GetString("mq.addr"),
        })
        if err != nil {
            return fmt.Errorf("MQ 连接失败: %w", err)
        }
        mqClient = client
        return nil
    },
    Close: func() error {
        if mqClient != nil {
            return mqClient.Close()
        }
        return nil
    },
})
```

## 5. 全局函数 vs 实例方法

### 场景 1：简单应用 - 使用全局函数

```go
func main() {
    crab.New()
    crab.UseLogx(...)
    crab.Run()
}
```

### 场景 2：测试或多实例 - 使用实例方法

```go
func TestApp(t *testing.T) {
    app := crab.New()
    app.Use(...)
    
    go func() {
        time.Sleep(time.Second)
        app.Done()  // 主动触发关闭
    }()
    
    if err := app.Run(); err != nil {
        t.Error(err)
    }
}
```

## 6. 生产环境配置建议

```go
func main() {
    app := crab.New()
    
    // 环境变量
    crab.UseEnv(env.WithAppKey("APP_ENV"))  // 从环境变量读取
    
    // 日志配置
    logLevel := logx.LevelInfo
    if env.IsProd() {
        logLevel = logx.LevelWarn  // 生产环境减少日志
    }
    
    crab.UseLogx(
        logx.WithLevel(logLevel),
        logx.WithEncodeJson(),      // 生产用 JSON
        logx.WithSource(true),      // 记录调用位置
    )
    
    // Gin 模式
    if env.IsProd() {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // 优雅关闭回调
    graceful.SetShutdownCallback(func(event string, data map[string]any) {
        logx.Info("关闭事件", "event", event, "data", data)
    })
    
    app.Run()
}
```

## 7. 常见陷阱

### ❌ 陷阱 1：在 Pre 阶段使用未初始化的资源

```go
// ❌ 错误：此时日志还未初始化
crab.Use(Handler{
    Pre: func() error {
        logx.Info("配置加载中")  // ❌ 如果 UseLogx 还没调用，这里会用默认配置
        return viperx.Build(...)
    },
})
crab.UseLogx(...)  // 日志在后面才初始化
```

### ❌ 陷阱 2：Init 返回 nil 但服务未启动

```go
// ❌ 错误：HTTP 服务未阻塞
crab.Use(Handler{
    Init: func() error {
        go server.ListenAndServe()  // ❌ goroutine 启动后立即返回
        return nil                   // Run() 会立即结束
    },
})

// ✅ 正确：HTTP 服务应该阻塞
crab.Use(Handler{
    Init: func() error {
        return server.ListenAndServe()  // ✅ 阻塞直到服务停止
    },
})
```

### ❌ 陷阱 3：Close 中处理顺序错误

```go
// ❌ 错误：先关数据库，再处理剩余请求
crab.Use(Handler{
    Close: func() error {
        db.Close()          // ❌ 数据库先关了
        return srv.Shutdown()  // 请求可能还在用数据库
    },
})

// ✅ 正确：先停止接收新请求，再关资源
crab.Use(Handler{
    Close: func() error {
        srv.Shutdown()      // ✅ 先停止服务
        return db.Close()   // 再关数据库
    },
})
```

## 8. 性能优化建议

### 使用连接池

```go
// Redis
crab.UseRedis(&redisx.Options{
    Addr:         "localhost:6379",
    PoolSize:     100,           // 连接池大小
    MinIdleConns: 10,            // 最小空闲连接
    MaxRetries:   3,             // 重试次数
}, &redisClient)

// MySQL
crab.UseMySQL(&mysqlx.ClientConfig{
    MaxIdleConns: 10,
    MaxOpenConns: 100,
    MaxLifetime:  time.Hour,
}, &mysqlClient)
```

### 合理设置超时

```go
// HTTP 服务器超时配置
server := &http.Server{
    Addr:           ":8080",
    Handler:        engine,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    IdleTimeout:    60 * time.Second,
    MaxHeaderBytes: 1 << 20,
}
crab.UseHTTP(server)
```

## 9. 监控和可观测性

```go
import (
    "github.com/bang-go/crab/core/pub/graceful"
)

func main() {
    app := crab.New()
    
    // 设置关闭回调，记录关闭事件
    graceful.SetShutdownCallback(func(event string, data map[string]any) {
        logx.Info("系统事件",
            "event", event,
            "timestamp", time.Now(),
            "data", data,
        )
    })
    
    // 健康检查端点
    crab.UseGin(":8080", func(engine *gin.Engine) {
        engine.GET("/health", func(c *gin.Context) {
            // 检查依赖服务健康状态
            c.JSON(200, gin.H{
                "status": "healthy",
                "env":    env.AppEnv(),
            })
        })
    })
    
    app.Run()
}
```

## 10. 单元测试

```go
func TestWorker(t *testing.T) {
    // 注意：测试中避免使用全局单例
    // New() 只会创建一次，后续测试会复用同一实例
    
    // 方案 1：单个测试文件内共享实例
    app := crab.New()
    
    // 方案 2：测试完整的应用启动流程
    go func() {
        time.Sleep(100 * time.Millisecond)
        app.Done()  // 主动触发关闭
    }()
    
    if err := app.Run(); err != nil {
        t.Fatal(err)
    }
}
```

## 总结

1. **初始化顺序**：环境 → 配置 → 日志 → 数据库 → HTTP
2. **错误处理**：框架不打日志，错误上抛到业务层
3. **Handler 选择**：Pre 用于前置条件，Init 用于资源初始化，Close 用于清理
4. **生产配置**：JSON 日志、Release 模式、合理的连接池和超时
5. **监控埋点**：健康检查、关闭回调、结构化日志
