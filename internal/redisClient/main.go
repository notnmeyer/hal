package redisClient

import (
	"github.com/go-redis/redis"
	"github.com/notnmeyer/hal/internal/config"
)

func New() *redis.Client {
	cfg := *make(config.Config).New()
	opts, err := redis.ParseURL(cfg["REDIS_URL"])
	if err != nil {
		panic(err.Error())
	}
	return redis.NewClient(opts)
}
