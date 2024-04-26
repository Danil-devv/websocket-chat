package kafka

import "github.com/IBM/sarama"

type Config struct {
	Brokers []string
	GroupID string
	Topics  []string
}

func InitConsumerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}
