package main

import (
	"fmt"
	"time"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
)

func main() {
	app := crab.New()
	defer app.Close() // 确保资源清理

	// Setup 阶段
	crab.Setup(
		func() error {
			return env.Build(env.WithAppEnv(env.DEV))
		},
		func() error {
			logx.Build(logx.WithLevel(logx.LevelInfo))
			return nil
		},
	)

	logx.Info("一次性 Job 示例", "env", env.AppEnv())

	// 一次性任务 - 只需要启动
	crab.Use(crab.StartOnly(func() error {
		logx.Info("开始处理数据")

		// 模拟数据处理
		for i := 1; i <= 5; i++ {
			fmt.Printf("处理进度: %d/5\n", i)
			time.Sleep(500 * time.Millisecond)
		}

		logx.Info("数据处理完成")
		return nil
	}))

	// Run - 执行任务
	if err := app.Run(); err != nil {
		logx.Error("任务执行失败", "error", err)
		return
	}

	// 任务执行完成，程序直接退出
	logx.Info("Job 完成")
}
