package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bang-go/crab"
)

// 模拟 Viper 的存储
var fakeViperData = make(map[string]string)

// 模拟 Viper.GetString
func ViperGetString(key string) string {
	val := fakeViperData[key]
	// 打印日志，方便观察是在什么时候调用的
	if val == "" {
		fmt.Printf("⚠️ [Viper] GetString(%q) 被调用 -> 返回空值 (配置未加载)\n", key)
	} else {
		fmt.Printf("✅ [Viper] GetString(%q) 被调用 -> 返回 %q\n", key, val)
	}
	return val
}

// 模拟构造函数 (典型的错误用法场景)
type Database struct {
	DSN string
}

func NewDatabase(dsn string) *Database {
	return &Database{DSN: dsn}
}

func main() {
	app := crab.New()

	// ----------------------------------------------------------------
	// 1. 注册配置加载 Hook
	// ----------------------------------------------------------------
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("\n--- [Step 1] 开始加载配置 ---")
			// 模拟加载配置
			fakeViperData["db.dsn"] = "mysql://root:123456@localhost:3306/mydb"
			fmt.Println("--- [Step 1] 配置加载完成 ---\n")
			return nil
		},
	})

	// ----------------------------------------------------------------
	// 2. 错误写法演示：在 Add 时就调用 ViperGetString
	// ----------------------------------------------------------------
	fmt.Println("--- [Setup] 准备注册错误的 Hook ---")
	// ❌ 错误：这里 ViperGetString 立即执行了！此时 Step 1 还没运行！
	wrongDSN := ViperGetString("db.dsn")

	// 即使我们把它传给构造函数，值也已经是空的了
	wrongDB := NewDatabase(wrongDSN)

	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("--- [Step 2] 启动错误的数据库连接 ---")
			fmt.Printf("❌ 错误连接: DSN = %q (预期是 mysql://...)\n", wrongDB.DSN)
			if wrongDB.DSN == "" {
				fmt.Println("   -> 失败原因: 在配置加载前就读取了配置值")
			}
			return nil
		},
	})

	// ----------------------------------------------------------------
	// 3. 正确写法演示：在 OnStart 内部调用 ViperGetString
	// ----------------------------------------------------------------
	fmt.Println("--- [Setup] 准备注册正确的 Hook ---")

	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("--- [Step 3] 启动正确的数据库连接 ---")

			// ✅ 正确：这里是在 Run() -> start() 流程中执行的
			// 此时 Step 1 已经执行完毕，配置已加载
			realDSN := ViperGetString("db.dsn")

			realDB := NewDatabase(realDSN)
			fmt.Printf("✅ 正确连接: DSN = %q\n", realDB.DSN)

			// 演示完成，退出
			go func() {
				time.Sleep(100 * time.Millisecond)
				app.Stop(context.Background())
			}()
			return nil
		},
	})

	// 运行应用
	fmt.Println("\n=== 应用启动 (crab.Run) ===")
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
