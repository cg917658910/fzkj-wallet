package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

func SetupRedis(ctx context.Context) {

	opt, err := redis.ParseURL("redis://root@localhost:6379/8")
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(opt)
}

func RedisClient() *redis.Client {
	return redisClient
}
