package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
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
	//receiveCh   chan *sarama.ConsumerMessage
	markCh <-chan *types.MarkMessageParams
}

func NewConsumerManager(ctx context.Context, ch chan *sarama.ConsumerMessage, markCh <-chan *types.MarkMessageParams) *MyConsumerManager {

	return &MyConsumerManager{
		consumer:    NewConsumerGroupHandler(ctx, ch),
		group:       NewConsumerGroup(),
		consumerNum: 10,
		markCh:      markCh,
		ctx:         ctx,
	}
}

func (m *MyConsumerManager) Start() error {
	logger.Info("Starting Consumer Manager...")
	if err := m.setupGroup(); err != nil {
		return err
	}
	if err := m.setupMarkChan(); err != nil {
		return err
	}
	if err := m.setupConsume(); err != nil {
		return err
	}
	//m.receiveMsg()
	logger.Info("Consumer Manager started successfully")

	return nil
}

func (m *MyConsumerManager) setupMarkChan() error {
	go func() {
		for {
			select {
			case msg, ok := <-m.markCh:
				if !ok {
					logger.Info("Mark channel closed")
					return
				}
				m.MarkMessage(msg)
			case <-m.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (m *MyConsumerManager) MarkMessage(msg *types.MarkMessageParams) {
	if m.consumer == nil {
		return
	}
	m.consumer.MarkMessage(msg.Msg, msg.MetaData)
}

func (m *MyConsumerManager) setupGroup() error {
	logger.Info("Starting consumer group...")
	if err := m.group.Setup(); err != nil {
		logger.Errorf("Failed to start consumer group: %v", err)
		return err
	}
	return nil
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
				logger.Info("Consumer group consumed message")
			}
		}()
	}
	return nil
}

func (m *MyConsumerManager) Stop() error {
	logger.Info("Stopping consumer group...")
	m.consumer.Commit()
	err := m.group.Cleanup()
	if err != nil {
		logger.Errorf("Failed to stop consumer group: %v", err)
		return err
	}
	logger.Info("Consumer group stopped successfully")
	return nil
}
