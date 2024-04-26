package config

import (
	"fmt"
	"os"
)

type Kafka struct {
	Brokers string
	Topics  string
	GroupID string
}

func getKafkaConfig() (Kafka, error) {
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
