package consumer

import (
	"context"

	"github.com/IBM/sarama"
	myConf "github.com/cg917658910/fzkj-wallet/notify-service/config"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
)

var logger = log.DLogger()

type ConsumerManager interface {
}
type MyConsumerManager struct {
	consumer    *consumerGroupHandler //kafka消费者
	group       *consumerGroup
	consumerNum int
	ctx         context.Context
	receiveCh   chan *sarama.ConsumerMessage
}

func NewConsumerManager(ctx context.Context) *MyConsumerManager {

	ch := make(chan *sarama.ConsumerMessage)
	return &MyConsumerManager{
		consumer:    NewConsumerGroupHandler(ch),
		group:       NewConsumerGroup(),
		consumerNum: 5,
		ctx:         ctx,
		receiveCh:   ch,
	}
}

func (m *MyConsumerManager) Start() error {
	logger.Info("Starting Consumer Manager...")
	if err := m.setupGroup(); err != nil {
		return err
	}
	if err := m.setupConsume(); err != nil {
		return err
	}
	m.receiveMsg()
	logger.Info("Consumer Manager started successfully")

	return nil
}

func (m *MyConsumerManager) setupGroup() error {
	logger.Info("Starting consumer group...")
	if err := m.group.Setup(); err != nil {
		logger.Errorf("Failed to start consumer group: %v", err)
		return err
	}
	return nil
}
func (m *MyConsumerManager) receiveMsg() {
	logger.Info("Starting to receive messages...")
	go func() {
		for {
			select {
			case msg := <-m.receiveCh:
				logger.Infof("Received message: %v", string(msg.Value))
				m.consumer.MarkMessage(msg, "") // 标记消息已消费
				// 处理消息
			case <-m.ctx.Done():
				logger.Info("Stopping message receiving...")
				return
			}
		}
	}()

}

func (m *MyConsumerManager) setupConsume() error {
	logger.Info("Starting to consume messages...")
	for range m.consumerNum {
		go func() {
			for {
				select {
				case <-m.ctx.Done():
					logger.Info("Stopping consumer...")
					return
				default:
				}
				err := m.group.group.Consume(m.ctx, []string{myConf.Configs.Kafka.OrderNofifyTopic}, m.consumer)
				if err != nil {
					logger.Errorf("Failed to consume message: %v", err)
					return
				}
			}
		}()
	}
	return nil
}

func (m *MyConsumerManager) Stop() error {
	logger.Info("Stopping consumer group...")
	err := m.group.Cleanup()
	if err != nil {
		logger.Errorf("Failed to stop consumer group: %v", err)
		return err
	}
	m.ctx.Done()
	logger.Info("Consumer group stopped successfully")
	return nil
}
