package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"storage/internal/adapters/kafka"
	"storage/internal/adapters/postgres"
	rds "storage/internal/adapters/redis"
)

type Config struct {
	Postgres *postgres.Config
	Kafka    *kafka.Config
	Redis    *rds.Config
}

func Get(logger *logrus.Logger, envFile string) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		logger.
			WithError(err).
			Error("cannot load .env file:")
	}

	postgresConfig, err := getPostgresConfig(logger)
	if err != nil {
		return nil, err
	}

	kafkaConfig, err := getKafkaConfig(logger)
	if err != nil {
		return nil, err
	}

	redisConfig, err := getRedisConfig(logger)
	if err != nil {
		return nil, err
	}

	config := &Config{
		Postgres: postgresConfig,
		Kafka:    kafkaConfig,
		Redis:    redisConfig,
	}
	return config, nil
}
