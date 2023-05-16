package redisx

import (
	"github.com/redis/go-redis/v9"
)

type Config redis.Options
type Client = redis.Client

func New(config Config) *Client {
	conf := redis.Options(config)
	client := redis.NewClient(&conf)
	return client
}
