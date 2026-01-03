package types

import "context"

// Runner 定义启动函数，接收 context 用于取消或超时控制
type Runner func(context.Context) error

// Stopper 定义停止函数，接收 context 用于超时控制
type Stopper func(context.Context) error
