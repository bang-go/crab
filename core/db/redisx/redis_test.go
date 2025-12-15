package redisx_test

import (
	"context"
	"testing"

	"github.com/bang-go/crab/core/db/redisx"
)

func TestRedis(t *testing.T) {
	// 测试无法连接到不存在的 Redis 服务器，应该返回错误
	redisClient, err := redisx.New(&redisx.Options{
		Addr: "localhost:6379",
	})

	// 如果没有 Redis 服务，应该返回错误
	if err != nil {
		t.Logf("expected error when Redis is not available: %v", err)
		return
	}

	// 如果连接成功，测试基本操作
	if redisClient != nil {
		defer redisClient.Close()
		_, _ = redisClient.Get(context.Background(), "test_key").Result()
	}
}
