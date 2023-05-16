package redisx

import (
	"github.com/redis/go-redis/v9"
)

type Config = redis.Options
type Client = redis.Client

var Nil = redis.Nil

func New(config Config) *Client {
	client := redis.NewClient(&config)
	return client
}
