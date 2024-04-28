package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"storage/internal/adapters/kafka"
	"strings"
)

type Kafka struct {
	Brokers string
	Topics  string
	GroupID string
}

func getKafkaConfig(logger logrus.FieldLogger) (*kafka.Config, error) {
	k, err := loadEnvKafkaConfig()
	if err != nil {
		return nil, err
	}
	kafkaConfig := &kafka.Config{
		Brokers: strings.Split(k.Brokers, ","),
		Topics:  strings.Split(k.Topics, ","),
		GroupID: k.GroupID,
		Logger:  logger.WithField("FROM", "[KAFKA-CONSUMER]"),
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
	groupID, ok := os.LookupEnv("KAFKA_GROUP_ID")
	if !ok {
		return Kafka{}, fmt.Errorf("KAFKA_GROUP_ID environment variable not set")
	}

	return Kafka{
		Brokers: brokers,
		Topics:  topics,
		GroupID: groupID,
	}, nil
}
