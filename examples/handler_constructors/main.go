package main

import (
	"fmt"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
)

func main() {
	app := crab.New()
	defer app.Close()

	// Setup 阶段 - 立即设置环境
	crab.Setup(
		func() error {
			return env.Build(env.WithAppEnv(env.DEV))
		},
		func() error {
			logx.Build(logx.WithLevel(logx.LevelInfo))
			return nil
		},
	)

	logx.Info("应用启动演示", "env", env.AppEnv())

	// 方式 1: StartClose - 需要启动和关闭
	crab.Use(crab.StartClose(
		func() error {
			logx.Info("启动数据库连接")
			// db.Connect()
			return nil
		},
		func() error {
			logx.Info("关闭数据库连接")
			// db.Close()
			return nil
		},
	))

	// 方式 2: StartOnly - 只需要启动
	crab.Use(
		crab.StartOnly(func() error {
			logx.Info("缓存预热")
			// cache.WarmUp()
			return nil
		}),
	)

	// 方式 3: 批量注册多个 StartOnly
	crab.Use(
		crab.StartOnly(func() error {
			logx.Info("加载初始数据")
			return nil
		}),
		crab.StartOnly(func() error {
			logx.Info("启动后台任务")
			go func() {
				// background task
			}()
			return nil
		}),
	)

	// 方式 4: 传统方式仍然可用（适合复杂场景）
	var server interface{ Shutdown() error }
	crab.Use(crab.Handler{
		Start: func() error {
			logx.Info("启动 HTTP 服务器")
			// server = startHTTPServer()
			return nil
		},
		Close: func() error {
			if server != nil {
				logx.Info("关闭 HTTP 服务器")
				return server.Shutdown()
			}
			return nil
		},
	})

	// 方式 5: CloseOnly - 只需要清理（极少见）
	var tempFile string
	crab.Setup(func() error {
		tempFile = "/tmp/app.tmp"
		fmt.Printf("创建临时文件: %s\n", tempFile)
		return nil
	})
	crab.Use(crab.CloseOnly(func() error {
		fmt.Printf("删除临时文件: %s\n", tempFile)
		return nil
	}))

	// 运行应用
	if err := app.Run(); err != nil {
		logx.Error("应用错误", "error", err)
	}
}
