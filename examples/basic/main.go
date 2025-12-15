package main

import (
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
)

func main() {
	app := crab.New()
	defer app.Close() // 确保资源清理

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
