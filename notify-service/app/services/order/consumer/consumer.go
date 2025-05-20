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

func (c *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	c.session = session
	for msg := range claim.Messages() {
		//logger.Infof("Received message: %s", string(msg.Value))
		c.sendMsgCh <- msg
	}
	return nil
}

// 标记已消费
func (c *consumerGroupHandler) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	c.session.MarkMessage(msg, metadata)
}
