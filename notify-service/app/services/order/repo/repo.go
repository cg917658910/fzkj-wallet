package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
	cacheLib "github.com/cg917658910/fzkj-wallet/notify-service/lib/cache"
	"github.com/redis/go-redis/v9"
)

type NotifyResultRepo struct {
	rdb                      *redis.Client
	notifyResultKeyPrefix    string
	notifyMsgOffsetKeyPrefix string
	ttl                      time.Duration
	ctx                      context.Context
	produceBFKey             string
}

var (
	_groupID = "cg.group.test.order.notify" //for test
	//_groupID = config.Configs.Kafka.OrderNofifyConsumerGroup
)

func NewNotifyResultRepo(ctx context.Context) *NotifyResultRepo {
	repo := &NotifyResultRepo{
		rdb:                      cacheLib.RedisClient(),
		notifyResultKeyPrefix:    "Cg:Notify:Order:Result",
		notifyMsgOffsetKeyPrefix: "Cg:Notify:Order:Offset",
		ctx:                      ctx,
		ttl:                      time.Hour * 1,
		produceBFKey:             "Cg:Notify:Order:Result:ProduceBF",
	}
	repo.init()

	return repo
}

func (repo *NotifyResultRepo) init() error {
	if repo.rdb == nil {
		return errors.New("init rdb is nil")
	}
	// init produce bf
	return repo.rdb.BFReserveWithArgs(repo.ctx, repo.produceBFKey, &redis.BFReserveOptions{
		Capacity: 10000 * 10,
		Error:    0.0001,
	}).Err()
}

func (repo *NotifyResultRepo) MarkNotifyMsgOffsetPending(topic string, partition int32, offset int64) error {
	key := markNotifyMsgOffsetPendingKey(repo.notifyMsgOffsetKeyPrefix, _groupID, topic, partition)
	return repo.rdb.ZAdd(repo.ctx, key, redis.Z{
		Score:  float64(offset),
		Member: offset,
	}).Err()
}

func (repo *NotifyResultRepo) MarkNotifyMsgOffsetDone(topic string, partition int32, offset int64) error {
	key := markNotifyMsgOffsetPendingKey(repo.notifyMsgOffsetKeyPrefix, _groupID, topic, partition)
	return repo.rdb.ZRem(repo.ctx, key, offset).Err()
}

func (repo *NotifyResultRepo) NotifyMsgPendingStatusMinOffset(topic string, partition int32, offset int64) (int64, error) {
	if repo.rdb == nil {
		return offset, errors.New("Notify Repo rdb is nil")
	}
	key := markNotifyMsgOffsetPendingKey(repo.notifyMsgOffsetKeyPrefix, _groupID, topic, partition)
	members, err := repo.rdb.ZRangeByScore(repo.ctx, key, &redis.ZRangeBy{Min: "-inf", Max: "+inf", Count: 1}).Result()
	if err != nil {
		return offset, err
	}
	fmt.Println("members: ", members)
	if len(members) == 0 {
		return offset, nil
	}
	member, _ := strconv.ParseInt(members[0], 10, 64)
	return min(offset, member), nil
}

func markNotifyMsgOffsetPendingKey(keyPrefix string, group string, topic string, partition int32) string {
	status := "pending"
	return fmt.Sprintf("%s:group_%s:topic_%s:status_%s:partition_%d", keyPrefix, _groupID, topic, status, partition)
}

func (repo *NotifyResultRepo) MarkProduceResult(id string) error {
	return repo.rdb.BFMAdd(repo.ctx, repo.produceBFKey, id).Err()
}

func (repo *NotifyResultRepo) ExistsProduceResutl(id string) (bool, error) {
	return repo.rdb.BFExists(repo.ctx, repo.produceBFKey, id).Result()
}

func (repo *NotifyResultRepo) Set(msg *types.NotifyResult) error {
	key := fmt.Sprintf("%s:%s", repo.notifyResultKeyPrefix, msg.MsgId)
	data, _ := json.Marshal(msg)
	return repo.rdb.Set(repo.ctx, key, data, repo.ttl).Err()
}

func (repo *NotifyResultRepo) Get(id string) (*types.NotifyResult, error) {

	key := fmt.Sprintf("%s:%s", repo.notifyResultKeyPrefix, id)
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
