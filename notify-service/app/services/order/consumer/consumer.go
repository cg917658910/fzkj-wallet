package consumer

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/repo"
)

// 消费者逻辑
type consumerGroupHandler struct {
	ctx       context.Context
	session   sarama.ConsumerGroupSession
	sendMsgCh chan *sarama.ConsumerMessage
	//*tracker.KafkaSafeConsumer
	notifyReop  *repo.NotifyResultRepo
	group       *consumerGroup
	needMarkNum uint64
}

func NewConsumerGroupHandler(ctx context.Context, group *consumerGroup, sendMsgCh chan *sarama.ConsumerMessage) *consumerGroupHandler {
	cgh := &consumerGroupHandler{
		ctx:         ctx,
		sendMsgCh:   sendMsgCh,
		group:       group,
		needMarkNum: 0,
		notifyReop:  repo.NewNotifyResultRepo(context.Background()),
	}
	atomic.StoreUint64(&cgh.needMarkNum, 0)
	return cgh
}
func (c *consumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	go func() {
		for {
			if c.canceled() {
				break
			}
			<-time.After(5 * time.Second)
			sess.Commit()
		}
	}()
	return nil
}

func (c *consumerGroupHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	// 清除前提交一次
	sess.Commit()
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.session = sess
	var (
		topic     = claim.Topic()
		partition = claim.Partition()
	)
	for msg := range claim.Messages() {
		logger.Infof("📩 消费者收到消息: %s", string(msg.Value))
		err := c.notifyReop.MarkNotifyMsgOffsetPending(topic, partition, msg.Offset)
		if err != nil {
			errLogger.Errorf("MarkNotifyMsgOffsetPending msg=%s, err= %s", string(msg.Value), err)
		}

		logger.Infof("needMarkNum=%d", atomic.LoadUint64(&c.needMarkNum))
		atomic.AddUint64(&c.needMarkNum, 1)

		c.sendMsgCh <- msg
	}
	return nil
}

func (c *consumerGroupHandler) canceled() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

func (c *consumerGroupHandler) Commit() {
	if c.session != nil {
		c.session.Commit()
	}
}

// 标记已消费
func (c *consumerGroupHandler) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	if msg == nil {
		return
	}
	logger.Infof("📩 消费者Mark消息: %s", string(msg.Value))
	if err := c.notifyReop.MarkNotifyMsgOffsetDone(msg.Topic, msg.Partition, msg.Offset); err != nil {
		errLogger.Errorf("MarkNotifyMsgOffsetDone msg=%s, err= %s", string(msg.Value), err)
	}
	c.session.MarkMessage(msg, metadata)
	atomic.AddUint64(&c.needMarkNum, ^uint64(0))

}
func (c *consumerGroupHandler) Stop(ctx context.Context) (err error) {
	logger.Infoln("Stop consumer group hanler...")
	defer func() {
		close(c.sendMsgCh)
		c.Commit() // 最後提交
		logger.Infoln("Stop consumer group hanler successfully")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		logger.Infof("Stop consumer group handler wait mark done|needMarkNum=%d", c.needMarkNum)
		if atomic.LoadUint64(&c.needMarkNum) == 0 {
			return
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
