package redisx

import (
	"github.com/redis/go-redis/v9"
)

type Config redis.Options

func New(config Config) *redis.Client {
	conf := redis.Options(config)
	client := redis.NewClient(&conf)
	return client
}
