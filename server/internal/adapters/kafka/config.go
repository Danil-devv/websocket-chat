package kafka

import "github.com/sirupsen/logrus"

type ProducerConfig struct {
	Brokers []string
	Topic   string
	Logger  *logrus.Logger
}
