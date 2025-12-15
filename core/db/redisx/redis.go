package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// 直接暴露第三方库的类型，方便业务方使用
type (
	Client  = redis.Client
	Options = redis.Options
)

// 常用常量
const (
	Nil = redis.Nil // Key 不存在错误
)

// DefaultOptions 返回默认配置
func DefaultOptions() *Options {
	return &Options{
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 5,
	}
}

// New 创建 Redis 客户端并验证连接
// 提供默认配置 + 连接验证，简化使用
func New(opts *Options) (*Client, error) {
	// 合并默认配置
	if opts.DialTimeout == 0 {
		opts.DialTimeout = 5 * time.Second
	}
	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = 3 * time.Second
	}
	if opts.WriteTimeout == 0 {
		opts.WriteTimeout = 3 * time.Second
	}
	if opts.PoolSize == 0 {
		opts.PoolSize = 10
	}

	client := redis.NewClient(opts)

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}

// MustNew 创建 Redis 客户端，失败则 panic
func MustNew(opts *Options) *Client {
	client, err := New(opts)
	if err != nil {
		panic(err)
	}
	return client
}
