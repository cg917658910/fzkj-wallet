package consumer

import (
	"context"
	"time"

	"github.com/IBM/sarama"
)

// 消费者逻辑
type consumerGroupHandler struct {
	ctx       context.Context
	session   sarama.ConsumerGroupSession
	sendMsgCh chan *sarama.ConsumerMessage
	//*tracker.KafkaSafeConsumer
}

func NewConsumerGroupHandler(ctx context.Context, sendMsgCh chan *sarama.ConsumerMessage) *consumerGroupHandler {
	return &consumerGroupHandler{
		ctx:       ctx,
		sendMsgCh: sendMsgCh,
	}
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

func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.session = session
	//tracker := c.GetTracker(claim.Partition())
	//tracker.SetInitOffset(claim.InitialOffset())
	for msg := range claim.Messages() {
		logger.Infof("📩 消费者收到消息: %s", string(msg.Value))
		if !c.canceled() {
			// tracker
			//tracker.Start(msg.Offset)
			c.sendMsgCh <- msg
		}
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
	//tracker := c.GetTracker(msg.Partition)
	//tracker.Done(msg.Offset)
	c.session.MarkMessage(msg, metadata)
	//Caller Notify Result url: http://localhost:8080/notifylogger.Infof("Consumer Manager Mark message: %s", string(msg.Value))
}
