package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bang-go/crab"
)

func main() {
	// 默认不配置 Logger
	app := crab.New()

	app.Add(crab.Hook{
		Name: "SilentComponent",
		OnStart: func(ctx context.Context) error {
			fmt.Println("组件业务逻辑执行中... (但不应该看到 Crab 的内部日志)")
			return nil
		},
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		app.Stop(context.Background())
	}()

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
