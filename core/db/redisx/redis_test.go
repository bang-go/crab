package redisx_test

import (
	"context"
	"github.com/bang-go/crab/core/db/redisx"
	"testing"
)

func TestRedis(t *testing.T) {
	var redisClient *redisx.Client
	redisClient = redisx.New(redisx.Config{})
	redisClient.Get(context.Background(), "")
}
