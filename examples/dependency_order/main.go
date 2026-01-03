package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bang-go/crab"
)

// 模拟一个全局配置管理器 (类似 viper)
var fakeViper = make(map[string]string)

func main() {
	app := crab.New()

	// 1. 注册配置初始化 Hook
	// 注意：Add 的顺序决定了 OnStart 的执行顺序
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("[Hook 1] 正在加载配置 (模拟 Viper)...")
			// 模拟加载耗时
			time.Sleep(100 * time.Millisecond)

			// 写入配置
			fakeViper["db.dsn"] = "postgres://user:pass@localhost:5432/mydb"
			fakeViper["app.name"] = "CrabDemo"

			fmt.Println("[Hook 1] 配置加载完成")
			return nil
		},
	})

	// 2. 注册依赖配置的 Hook
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("[Hook 2] 正在初始化数据库连接...")

			// 尝试获取配置
			dsn := fakeViper["db.dsn"]
			if dsn == "" {
				return fmt.Errorf("关键配置 db.dsn 未找到！配置加载失败？")
			}

			fmt.Printf("[Hook 2] 成功读取配置 db.dsn: %s\n", dsn)
			fmt.Println("[Hook 2] 数据库连接成功")
			return nil
		},
	})

	// 3. 注册业务逻辑 Hook
	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			appName := fakeViper["app.name"]
			fmt.Printf("[Hook 3] 应用 %s 已就绪\n", appName)

			// 演示完成，自动退出
			go func() {
				time.Sleep(100 * time.Millisecond)
				app.Stop(context.Background())
			}()
			return nil
		},
	})

	if err := app.Run(); err != nil {
		log.Fatalf("应用启动失败: %v", err)
	}
}
