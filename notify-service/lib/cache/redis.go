package cache

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/notify-service/config"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
)

func SetupRedis(ctx context.Context) {

	opt, err := redis.ParseURL(config.Configs.Redis.URL)
	if err != nil {
		panic(err)
	}
	redisClient = redis.NewClient(opt)
}

func RedisClient() *redis.Client {
	return redisClient
}
