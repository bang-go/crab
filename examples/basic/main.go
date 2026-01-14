package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bang-go/crab"
)

func main() {
	app := crab.New()

	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Setting up environment...")
			return nil
		},
	})

	app.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("应用启动")
			return nil
		},
	})

	app.OnShutdown(func() {
		log.Println("应用正在关闭...")
	})

	if err := app.Run(); err != nil {
		log.Printf("应用错误: %v", err)
	}
}
