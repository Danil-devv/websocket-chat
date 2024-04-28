package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"server/internal/adapters/kafka"
	"strings"
)

type Kafka struct {
	Brokers string
	Topic   string
}

func getKafkaConfig(logger logrus.FieldLogger) (*kafka.Config, error) {
	k, err := loadEnvKafkaConfig()
	if err != nil {
		return nil, err
	}
	kafkaConfig := &kafka.Config{
		Brokers: strings.Split(k.Brokers, ","),
		Topic:   k.Topic,
		Logger:  logger.WithField("FROM", "[KAFKA-PRODUCER]"),
	}
	return kafkaConfig, nil
}

func loadEnvKafkaConfig() (Kafka, error) {
	brokers, ok := os.LookupEnv("KAFKA_BROKERS")
	if !ok {
		return Kafka{}, fmt.Errorf("KAFKA_BROKERS environment variable not set")
	}
	topics, ok := os.LookupEnv("KAFKA_TOPICS")
	if !ok {
		return Kafka{}, fmt.Errorf("KAFKA_TOPICS environment variable not set")
	}

	return Kafka{
		Brokers: brokers,
		Topic:   topics,
	}, nil
}
