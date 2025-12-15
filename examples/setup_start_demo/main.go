package main

import (
	"fmt"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/base/viperx"
)

func main() {
	app := crab.New()
	defer app.Close()

	// Setup 阶段 - 立即设置环境（准备阶段）
	crab.Setup(
		func() error {
			fmt.Println("1. 设置环境变量")
			return env.Build(env.WithAppEnv(env.DEV))
		},
		func() error {
			fmt.Println("2. 加载配置文件")
			return viperx.Build(&viperx.Config{
				ConfigFormat: viperx.FileFormatYaml,
				ConfigPaths:  []string{".", "./config"},
				ConfigNames:  []string{"config"},
			})
		},
		func() error {
			fmt.Println("3. 初始化日志")
			logx.Build(logx.WithLevel(logx.LevelInfo))
			return nil
		},
	)

	logx.Info("环境准备完成", "env", env.AppEnv())

	// Use 阶段 - 注册资源生命周期
	crab.Use(crab.Handler{
		Start: func() error {
			logx.Info("启动定时任务")
			return nil
		},
		Close: func() error {
			logx.Info("停止定时任务")
			return nil
		},
	})

	// Run 阶段 - 运行应用（触发所有 Start）
	logx.Info("应用启动")
	if err := app.Run(); err != nil {
		logx.Error("应用错误", "error", err)
	}
}
