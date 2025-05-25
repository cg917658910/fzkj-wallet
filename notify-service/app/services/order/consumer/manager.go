package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
	myConf "github.com/cg917658910/fzkj-wallet/notify-service/config"
)

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

	group := NewConsumerGroup()
	return &MyConsumerManager{
		consumer:    NewConsumerGroupHandler(ctx, group, ch),
		group:       group,
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
	logger.Info("Consumer Manager started successfully")

	return nil
}

func (m *MyConsumerManager) setupMarkChan() error {
	go func() {
		for msg := range m.markCh {
			m.MarkMessage(msg)
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
				}
				logger.Info("Consumer group consumed message")
			}
		}()
	}
	return nil
}

func (m *MyConsumerManager) Stop(ctx context.Context) error {
	logger.Info("Stopping consumer manager...")
	//pause all stop consumer messsage
	if err := m.group.PrepareStop(); err != nil {
		logger.Errorf("Failed to PrepareStop consumer group: %v", err)
		return err
	}
	// wait mark done
	if err := m.consumer.Stop(ctx); err != nil {
		logger.Errorf("Failed to stop consumer handler: %v", err)
		return err
	}
	//close group
	if err := m.group.Stop(); err != nil {
		logger.Errorf("Failed to Stop consumer group: %v", err)
		return err
	}
	logger.Info("Stopping consumer manager successfully")
	return nil
}
