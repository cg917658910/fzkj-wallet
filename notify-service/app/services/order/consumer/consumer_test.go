package consumer

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/config"
)

func init() {

}

func TestConsumerGroupHandlerStop(t *testing.T) {
	gp := NewConsumerGroup()
	ch := make(chan *sarama.ConsumerMessage)
	c := NewConsumerGroupHandler(context.Background(), gp, ch)

	c.Stop(context.Background())
}

type myConsumer struct {
	needMarkNum uint64
}

func NewMy() *myConsumer {
	return &myConsumer{}
}
func TestConsumerGroupHandlerNeedMarkNum(t *testing.T) {
	/* gp := NewConsumerGroup()
	ch := make(chan *sarama.ConsumerMessage)
	c := NewConsumerGroupHandler(context.Background(), gp, ch)
	var n uint64 */
	c := NewMy()
	atomic.AddUint64(&c.needMarkNum, 10)
	//c.incMarkNum()
	num := atomic.LoadUint64(&c.needMarkNum)
	if num != 1 {
		t.Errorf("expected  1, got %d", num)
	}
}
func TestConsumerLog(t *testing.T) {
	config.Setup()
	logger.Infoln("test log")
}
