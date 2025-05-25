package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/repo"
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
	notifyReop     *repo.NotifyResultRepo
}

const ()

func NewProducerManager(ctx context.Context, notifyResultCh <-chan *types.NotifyResult, markCh chan *types.MarkMessageParams) *myProducerManager {
	return &myProducerManager{
		ctx:            ctx,
		notifyResultCh: notifyResultCh,
		workerNum:      500,
		markMessageCh:  markCh,
		notifyReop:     repo.NewNotifyResultRepo(ctx),
	}
}

func (m *myProducerManager) Start() error {
	// 启动生产者
	logger.Info("Starting Producer Manager...")
	if err := m.setupProducer(); err != nil {
		errLogger.Errorf("Failed to setup producer: %v", err)
		return err
	}
	if err := m.setupWorker(); err != nil {
		errLogger.Errorf("Failed to setup producer: %v", err)
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
			for msg := range m.notifyResultCh {
				logger.Infof("producer received msg=%s", msg.MsgId)
				m.processWorker(msg)
			}
			/* for {
				select {
				case msg := <-m.notifyResultCh:
					m.processWorker(msg)
				case <-m.ctx.Done():
					return
				}
			} */
		}()
	}
	return nil
}

func (m *myProducerManager) processWorker(msg *types.NotifyResult) error {
	if msg == nil {
		return nil
	}
	var isSend bool
	// 1.检查消息是否已经发送
	if m.notifyReop != nil && msg.MsgId != "" {
		//isSend, _ = m.notifyReop.ExistsProduceResutl(msg.MsgId)
	}
	if isSend {
		return m.markMessage(msg.RawMsg, msg.Msg)
	}
	// 发送消息
	if err := m.produceMessage(msg); err != nil {
		errLogger.Errorf("Producer produceMessage failed: %s", err)
		return err
	}
	// mark 原消息
	m.markMessage(msg.RawMsg, msg.Msg)
	return nil
}

func (m *myProducerManager) markMessage(msg *sarama.ConsumerMessage, metadata string) error {
	// TODO: 通道安全检测
	if m.markMessageCh == nil {
		logger.Warnln("Producer Manager markMessageCh is nil")
		return nil
	}
	m.markMessageCh <- &types.MarkMessageParams{
		Msg:      msg,
		MetaData: metadata,
	}
	if !m.canceled() {

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

func (m *myProducerManager) produceMessage(data *types.NotifyResult) (err error) {
	if m.producer == nil {
		return genError("Producer Manager kafka producer is nil")
	}
	if data == nil {
		return genError("message is nil")
	}
	key := fmt.Sprintf("%s", data.MsgId)
	valueByte, err := json.Marshal(data)
	if err != nil {
		err = genError(fmt.Sprintf("❌生产者发送消息失败|tag=消息序列化失败|key=%s|error=%v", key, err))
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: getTopicNameByPlatform(data.Platform),
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(valueByte),
	}
	partition, offset, err := m.producer.SendMessage(msg)
	if err != nil {
		err = genError(fmt.Sprintf("❌ 生产者发送消息失败|tag=producerSend|key=%s|error=%v", key, err))
		return
	}
	if err = m.notifyReop.MarkProduceResult(data.MsgId); err != nil {
		err = genError(fmt.Sprintf("❌ 生产者发送消息失败|tag=MarkProduceResult失败|key=%s|error=%v", key, err))
		return
	}
	logger.Infof("✅ 生产者发送消息: |key=%s|Partition=%d|Offset=%d|", key, partition, offset)

	return nil
}

func getTopicNameByPlatform(platform string) string {
	topics := strings.Split(config.Configs.Kafka.OrderNofifyResultTopics, ",")
	for _, topic := range topics {
		if strings.HasPrefix(topic, platform+".") {
			return topic
		}
	}
	return config.Configs.Kafka.OrderNofifyResultDefaultTopic
}

func (m *myProducerManager) Stop() error {
	logger.Info("Stopping Producer Manager...")
	if m.producer != nil {
		// TODO: 关闭时机
		if err := m.producer.Close(); err != nil {
			errLogger.Errorf("Failed to close producer: %v", err)
			return err
		}
		m.producer = nil
		logger.Info("Producer closed successfully")
	}
	if err := m.closeMarkCh(); err != nil {
		errLogger.Errorf("Failed to close mark channel: %v", err)
		return err
	}
	logger.Info("Producer Manager stopped successfully")
	return nil
}

func (m *myProducerManager) closeMarkCh() error {
	logger.Infoln("Producer Manager close mark channel...")
	sync.OnceFunc(func() {
		close(m.markMessageCh)
	})
	logger.Infoln("Producer Manager close mark channel successfully")
	return nil
}
