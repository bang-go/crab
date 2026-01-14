# Crab 框架架构指南

本文档定义了使用 `crab` 框架的标准模式，重点介绍用于依赖注入（DI）场景下的 **生命周期注入模式 (Lifecycle Injection Pattern)**。

## 1. 核心理念：生命周期注入

管理组件生命周期（连接池、后台工作者、长连接等）的推荐模式是 **生命周期注入**。

在这种模式下，组件提供者（Provider）不应该返回清理函数或 `Hook` 对象，而是通过接收 `crab.Lifecycle` 接口并自行注册生命周期行为。

### 核心优势
*   **消除样板代码 (Zero Boilerplate)**：不需要定义复杂的中间结构体来承接清理函数，也不需要 DI 工具进行字段提取。
*   **高内聚 (Encapsulation)**：资源的创建逻辑与其销毁逻辑封装在同一个函数内。
*   **接口导向 (Interface Driven)**：Provider 依赖于抽象的 `Lifecycle` 接口，而非具体的框架实现。

## 2. 标准实施模式

### 步骤 1: 资源提供者 (Provider Layer)

Provider 函数声明对 `crab.Lifecycle` 的依赖，并在内部注册钩子。

**示例：**
```go
import "github.com/bang-go/crab"

func NewResource(conf *Config, lc crab.Lifecycle) (*Resource, error) {
    res := &Resource{} // 初始化资源
    
    // 自动注册生命周期钩子
    // 方式 A：使用辅助函数（最简）
    lc.Append(crab.Close(res.Close))
    
    // 方式 B：使用完整 Hook（当需要自定义启动/停止逻辑或名称时）
    lc.Append(crab.Hook{
        Name: "my-resource",
        OnStart: func(ctx context.Context) error {
            return res.Open(ctx)
        },
        OnStop: func(ctx context.Context) error {
            return res.Close()
        },
    })
    
    return res, nil
}
```

### 步骤 2: 依赖注入配置 (DI)

在 DI 容器配置中，需要将 `*crab.Registry` 绑定到 `crab.Lifecycle` 接口。

```go
// 以常见的 DI 逻辑为例
var ProviderSet = NewSet(
    // 1. 提供 Registry 单例
    crab.NewRegistry, 
    // 2. 将 Registry 绑定到接口
    Bind(new(crab.Lifecycle), new(*crab.Registry)),
    
    // ... 其他 Provider ...
)
```

### 步骤 3: 应用启动 (Main)

`main` 函数负责获取收集满钩子的 `Registry`，并将其注入到 `App` 实例中执行。

```go
func main() {
    // 初始化 DI 容器，获取 registry
    container, registry, err := initContainer()
    if err != nil {
        panic(err)
    }

    app := crab.New()
    app.Add(registry.Hooks()...)

    // 将 registry 收集的所有钩子一次性注入并运行
    if err := app.Run(); err != nil {
        panic(err)
    }
}
```

## 3. 核心 API

### 生命周期管理
*   **`crab.Lifecycle` (接口)**：定义了 `Append(Hook)` 方法。
*   **`crab.Registry` (结构体)**：`Lifecycle` 的线程安全实现，用于收集钩子。

### 快捷构造器 (Helpers)
*   **`crab.Close(fn)`**：将不带参数的 `Close() error` 函数转换为 Hook。
*   **`crab.OnStart(fn)`**：仅定义启动逻辑。
*   **`crab.OnStop(fn)`**：仅定义停止逻辑。

## 4. 最佳实践原则

1.  **就近注册**：谁创建资源，谁负责调用 `lc.Append` 注册销毁逻辑。
2.  **避免返回 Cleanup**：函数签名应尽可能保持为 `func(...) (*T, error)`。
3.  **命名规范**：在手动创建 `Hook` 时，务必提供 `Name` 字段，以便在日志中追踪各组件的执行耗时。

---
*由 Gemini Agent 生成，作为项目通用架构规范。*
