package consumer

import (
	"strings"

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
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                            // 从最早的消息开始消费
	config.Consumer.Offsets.AutoCommit.Enable = false                                //关闭自动提交
	config.Consumer.Return.Errors = true
	/* config.Consumer.Group.Session.Timeout = time.Minute * 10   // 10分钟
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Minute // 3分钟 */
	/* config.Consumer.Group.Rebalance.Timeout = 10 * time.Minute // 10分钟
	config.Consumer.Group.Rebalance.Retry.Max = 10
	config.Consumer.Group.Rebalance.Retry.Backoff = 10 * time.Second // 10秒 */

	brokers := strings.Split(myConf.Configs.Kafka.Brokers, ",")

	group, err := sarama.NewConsumerGroup(brokers, myConf.Configs.Kafka.OrderNofifyConsumerGroup, config)
	if err != nil {
		logger.Fatalf("创建消费者组失败: %v", err)
	}
	c.group = group
	// 监听 Kafka 消费者组错误
	go func() {
		for err := range group.Errors() {
			errLogger.Errorf("消费者组错误: %v\n", err)
		}
	}()

	return nil
}

func (c *consumerGroup) PrepareStop() error {
	logger.Info("PauseAll consumer group...")
	c.group.PauseAll()
	logger.Info("Consumer group PauseAll successfully")
	return nil
}

func (c *consumerGroup) Stop() error {
	logger.Info("Stopping consumer group...")
	if err := c.group.Close(); err != nil {
		errLogger.Errorf("关闭消费者组失败: %v", err)
		return err
	}
	logger.Info("Consumer group stopped successfully")
	return nil
}
