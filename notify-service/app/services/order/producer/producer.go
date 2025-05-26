package producer

import (
	"strings"

	"github.com/IBM/sarama"
	myConf "github.com/cg917658910/fzkj-wallet/notify-service/config"
)

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	//config.Producer.Idempotent = true
	config.Producer.Return.Successes = true
	//config.Net.MaxOpenRequests = 1
	//config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认 TODO: 需要权衡
	config.Producer.Retry.Max = 5
	config.Producer.Partitioner = sarama.NewHashPartitioner // 通过 Key 进行分区
	brokers := strings.Split(myConf.Configs.Kafka.Brokers, ",")
	return sarama.NewSyncProducer(brokers, config)
}
