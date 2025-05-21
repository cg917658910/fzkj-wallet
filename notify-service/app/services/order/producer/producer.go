package producer

import (
	"strings"

	"github.com/IBM/sarama"
	myConf "github.com/cg917658910/fzkj-wallet/notify-service/config"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
)

var logger = log.DLogger()

func newProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner // 通过 Key 进行分区
	brokers := strings.Split(myConf.Configs.Kafka.Brokers, ",")
	return sarama.NewSyncProducer(brokers, config)
}
