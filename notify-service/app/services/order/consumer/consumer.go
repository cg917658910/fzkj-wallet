package consumer

import (
	"time"

	"github.com/IBM/sarama"
)

// 消费者逻辑
type consumerGroupHandler struct {
	session   sarama.ConsumerGroupSession
	sendMsgCh chan *sarama.ConsumerMessage
}

func NewConsumerGroupHandler(sendMsgCh chan *sarama.ConsumerMessage) *consumerGroupHandler {
	return &consumerGroupHandler{
		sendMsgCh: sendMsgCh,
	}
}
func (c *consumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	go func() {
		for {
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
	for msg := range claim.Messages() {
		logger.Infof("📩 消费者收到消息: %s", string(msg.Value))
		c.sendMsgCh <- msg
	}
	return nil
}

// 标记已消费
func (c *consumerGroupHandler) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	c.session.MarkMessage(msg, metadata)
	logger.Infof("Consumer Manager Mark message: %s", string(msg.Value))
}
