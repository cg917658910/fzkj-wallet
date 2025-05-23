package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
	cacheLib "github.com/cg917658910/fzkj-wallet/notify-service/lib/cache"
	"github.com/redis/go-redis/v9"
)

type NotifyResultRepo struct {
	rdb       *redis.Client
	keyPrefix string
	ttl       time.Duration
	ctx       context.Context
}

func NewNotifyResultRepo(ctx context.Context) *NotifyResultRepo {

	return &NotifyResultRepo{

		rdb:       cacheLib.RedisClient(),
		keyPrefix: "Cg:Notify:Order:Result",
		ctx:       ctx,
		ttl:       time.Hour * 1,
	}
}

func (repo *NotifyResultRepo) Set(msg *types.NotifyResult) error {
	key := fmt.Sprintf("%s:%s", repo.keyPrefix, msg.MsgId)
	data, _ := json.Marshal(msg)
	return repo.rdb.Set(repo.ctx, key, data, repo.ttl).Err()
}

func (repo *NotifyResultRepo) Get(id string) (*types.NotifyResult, error) {

	key := fmt.Sprintf("%s:%s", repo.keyPrefix, id)
	val, err := repo.rdb.Get(repo.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var u types.NotifyResult
	err = json.Unmarshal([]byte(val), &u)
	return &u, err
}

func (repo *NotifyResultRepo) Del(key string) (err error) {
	return repo.rdb.Del(repo.ctx, key).Err()
}
