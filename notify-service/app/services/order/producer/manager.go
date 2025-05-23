package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/types"
	"github.com/cg917658910/fzkj-wallet/notify-service/config"
)

// Producer 代表生产者的接口类型。
type ProducerManager interface {
	// Start 用于启动生产者。
	Start() error
	// Stop 用于停止生产者。
	Stop() error
}

type myProducerManager struct {
	ctx            context.Context
	notifyResultCh <-chan *types.NotifyResult
	workerNum      int
	producer       sarama.SyncProducer
	markMessageCh  chan *types.MarkMessageParams // 用于标记消息的通道
}

const ()

func NewProducerManager(ctx context.Context, notifyResultCh <-chan *types.NotifyResult, markCh chan *types.MarkMessageParams) *myProducerManager {
	return &myProducerManager{
		ctx:            ctx,
		notifyResultCh: notifyResultCh,
		workerNum:      100,
		markMessageCh:  markCh,
	}
}

func (m *myProducerManager) Start() error {
	// 启动生产者
	logger.Info("Starting Producer Manager...")
	if err := m.setupProducer(); err != nil {
		logger.Errorf("Failed to setup producer: %v", err)
		return err
	}
	if err := m.setupWorker(); err != nil {
		logger.Errorf("Failed to setup producer: %v", err)
		return err
	}
	return nil
}

func (m *myProducerManager) setupProducer() error {
	producer, err := newProducer()
	if err != nil {
		return err
	}
	m.producer = producer
	return nil
}

func (m *myProducerManager) setupWorker() error {
	for range m.workerNum {
		go func() {
			for {
				select {
				case msg := <-m.notifyResultCh:
					m.processWorker(msg)
				case <-m.ctx.Done():
					return
				}
			}
		}()
	}
	return nil
}

func (m *myProducerManager) processWorker(msg *types.NotifyResult) error {
	// 1.发送消息
	if err := m.produceMessage(msg); err != nil {
		return err
	}
	// 2. mark 原消息
	m.markMessage(msg.RawMsg, msg.Msg)
	return nil
}

func (m *myProducerManager) markMessage(msg *sarama.ConsumerMessage, metadata string) error {
	// TODO: 通道安全检测
	if m.markMessageCh == nil {
		logger.Warnln("Producer Manager markMessageCh is nil")
		return nil
	}
	if !m.canceled() {
		m.markMessageCh <- &types.MarkMessageParams{
			Msg:      msg,
			MetaData: metadata,
		}
	}

	return nil
}

func (m *myProducerManager) canceled() bool {
	select {
	case <-m.ctx.Done():
		return true
	default:
		return false
	}
}

func (m *myProducerManager) produceMessage(data *types.NotifyResult) error {
	if m.producer == nil {
		return genError("Producer Manager kafka producer is nil")
	}
	if data == nil {
		return genError("message is nil")
	}
	key := fmt.Sprintf("%s:%s", data.Platform, data.Data.DataId)
	valueByte, err := json.Marshal(data)
	if err != nil {
		logger.Warnf("❌ 生产者 %s 消息序列化失败: %v", key, err)
	}
	msg := &sarama.ProducerMessage{
		Topic: getTopicNameByPlatform(data.Platform),
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(valueByte),
	}

	if m.producer == nil {
		logger.Warnf("❌ 生产者 Prodcuer is nil ")
		return nil
	}
	if !m.canceled() {
		_, _, err = m.producer.SendMessage(msg)
		if err != nil {
			logger.Errorf("❌ 生产者 %s 发送消息失败: %v", key, err)
			return err
		}
	}

	//logger.Infof("✅ 生产者发送消息: %s (Partition=%d, Offset=%d)", key, partition, offset)
	return nil
}

func getTopicNameByPlatform(platform string) string {
	topics := strings.SplitSeq(config.Configs.Kafka.OrderNofifyResultTopics, ",")
	for topic := range topics {
		if strings.HasPrefix(topic, platform+".") {
			return topic
		}
	}

	return config.Configs.Kafka.OrderNofifyResultDefaultTopic
}

func (m *myProducerManager) Stop() error {
	logger.Info("Stopping Producer Manager...")
	if m.producer != nil {
		if err := m.producer.Close(); err != nil {
			logger.Errorf("Failed to close producer: %v", err)
			return err
		}
		m.producer = nil
		logger.Info("Producer closed successfully")
	}
	logger.Info("Producer Manager stopped successfully")
	return nil
}
