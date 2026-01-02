package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bang-go/crab"
)

func main() {
	// 方式 1: 使用默认实例
	crab.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Setting up environment...")
			return nil
		},
	})

	crab.Add(crab.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("应用启动")
			return nil
		},
	})

	if err := crab.Run(); err != nil {
		log.Printf("应用错误: %v", err)
	}
}
