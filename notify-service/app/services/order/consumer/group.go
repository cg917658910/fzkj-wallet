package consumer

import (
	"time"

	"github.com/IBM/sarama"
	myConf "github.com/cg917658910/fzkj-wallet/notify-service/config"
)

type consumerGroup struct {
	group sarama.ConsumerGroup
}

func NewConsumerGroup() *consumerGroup {
	return &consumerGroup{}
}

func (c *consumerGroup) Setup() error {
	config := sarama.NewConfig()
	config.Version = sarama.V3_2_3_0

	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin() // 轮询分区
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = false //关闭自动提交
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Session.Timeout = time.Minute * 10 // 10分钟
	/* config.Consumer.Group.Heartbeat.Interval = 3 * 60 * 1000 // 3分钟
	config.Consumer.Group.Rebalance.Timeout = 10 * 60 * 1000 // 10分钟
	config.Consumer.Group.Rebalance.Retry.Max = 10
	config.Consumer.Group.Rebalance.Retry.Backoff = 10 * 1000 // 10秒 */

	group, err := sarama.NewConsumerGroup([]string{myConf.Configs.Kafka.Brokers}, myConf.Configs.Kafka.OrderNofifyConsumerGroup, config)
	if err != nil {
		logger.Fatalf("创建消费者组失败: %v", err)
	}
	c.group = group
	// 监听 Kafka 消费者组错误
	go func() {
		for err := range group.Errors() {
			logger.Errorf("消费者组错误: %v\n", err)
		}
	}()

	return nil
}

func (c *consumerGroup) Cleanup() error {
	if err := c.group.Close(); err != nil {
		logger.Errorf("关闭消费者组失败: %v", err)
		return err
	}
	return nil
}
