package config

import (
	"chat/internal/adapters/kafka"
	"chat/internal/adapters/postgres"
	rds "chat/internal/adapters/redis"
	"chat/internal/adapters/websocket"
	"chat/internal/app"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Postgres *postgres.Config
	Kafka    *kafka.Config
	Redis    *rds.Config
	Server   *websocket.Config
	App      *app.Config
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

	serverConfig, err := getServerConfig()
	if err != nil {
		return nil, err
	}

	appConfig, err := getAppConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{
		Postgres: postgresConfig,
		Kafka:    kafkaConfig,
		Redis:    redisConfig,
		Server:   serverConfig,
		App:      appConfig,
	}
	return config, nil
}
