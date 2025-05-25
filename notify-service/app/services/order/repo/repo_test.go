package repo

import (
	"context"
	"strings"
	"testing"

	"github.com/cg917658910/fzkj-wallet/notify-service/lib/cache"
)

func TestInitCreateProduceBF(t *testing.T) {
	ctx := context.Background()
	cache.SetupRedis(ctx)
	ctl := NewNotifyResultRepo(ctx)
	err := ctl.init()
	var errItemExists = "ERR item exists"
	if !strings.Contains(err.Error(), errItemExists) {
		t.Errorf("expected %s, got %s", errItemExists, err)
	}
}

func TestMarkProduceResult(t *testing.T) {
	ctx := context.Background()
	cache.SetupRedis(ctx)
	ctl := NewNotifyResultRepo(ctx)

	id1 := "msg_id_1"
	err := ctl.MarkProduceResult(id1)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	isExist, err := ctl.ExistsProduceResutl(id1)
	if err != nil {
		t.Errorf("exist expected nil, got %s", err)
	}
	if !isExist {
		t.Errorf("expected true, got %v", isExist)
	}
}

func TestMarkNotifyMsgOffsetPending(t *testing.T) {
	ctx := context.Background()
	cache.SetupRedis(ctx)
	ctl := NewNotifyResultRepo(ctx)
	var (
		topic1           = "topic1"
		topic2           = "topic2"
		partition0 int32 = 0
		partition1 int32 = 1
		partition2 int32 = 2
	)

	ctl.MarkNotifyMsgOffsetPending(topic1, partition0, 100)
	ctl.MarkNotifyMsgOffsetPending(topic1, partition0, 99)
	ctl.MarkNotifyMsgOffsetPending(topic1, partition0, 60)
	ctl.MarkNotifyMsgOffsetPending(topic1, partition0, 30)
	ctl.MarkNotifyMsgOffsetPending(topic1, partition0, 11)

	ctl.MarkNotifyMsgOffsetPending(topic1, partition1, 100)
	err := ctl.MarkNotifyMsgOffsetPending(topic2, partition2, 1)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	ctl.MarkNotifyMsgOffsetDone(topic1, partition0, 11)
	err = ctl.MarkNotifyMsgOffsetDone(topic1, partition1, 1000)
	if err != nil {
		t.Errorf("done offset 1000, expected nil, got %s", err)
	}
	minPendingOffset, err := ctl.NotifyMsgPendingStatusMinOffset(topic1, partition0, 200)
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
	if minPendingOffset != 30 {
		t.Errorf("expected 30, got %d", minPendingOffset)
	}
	minPendingOffset, err = ctl.NotifyMsgPendingStatusMinOffset(topic1, partition0, 20)
	if minPendingOffset != 20 {
		t.Errorf("expected 20, got %d", minPendingOffset)
	}
}
