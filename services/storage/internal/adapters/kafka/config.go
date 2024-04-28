package kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Brokers []string
	GroupID string
	Topics  []string
	Logger  logrus.FieldLogger
}

func InitConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}
