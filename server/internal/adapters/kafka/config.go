package kafka

import "github.com/sirupsen/logrus"

type Config struct {
	Brokers []string
	Topic   string
	Logger  logrus.FieldLogger
}
