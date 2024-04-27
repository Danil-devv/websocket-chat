package config

import (
	"context"
	"fmt"
	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"storage/internal/adapters/kafka"
	"storage/internal/adapters/postgres"
	rds "storage/internal/adapters/redis"
	"strings"
)

type config struct {
	Kafka
	Postgres
	Redis
}

func Get(logger *logrus.Logger, envFile string) (*postgres.Config, *kafka.Config, *rds.Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		logger.
			WithError(err).
			Error("cannot load .env file:")
	}

	cfg, err := parse()
	if err != nil {
		logger.
			WithError(err).
			Error("cannot get config")
		return nil, nil, nil, err
	}

	pgxConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			cfg.Postgres.Username,
			cfg.Postgres.Password,
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Database,
		),
	)
	if err != nil {
		logger.
			WithError(err).
			Error("cannot parse postgres config")
		return nil, nil, nil, err
	}
	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(logger.WithField("FROM", "[PGX-POOL]")),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		logger.
			WithError(err).
			Error("cannot create pgx pool")
		return nil, nil, nil, err
	}
	postgresConfig := &postgres.Config{
		Pool:   pool,
		Logger: logger.WithField("FROM", "[POSTGRES]"),
	}

	kafkaConfig := &kafka.Config{
		Brokers: strings.Split(cfg.Kafka.Brokers, ","),
		Topics:  strings.Split(cfg.Kafka.Topics, ","),
		GroupID: cfg.Kafka.GroupID,
		Logger:  logger.WithField("FROM", "[KAFKA-CONSUMER]"),
	}
	redisConfig := &rds.Config{
		Opt: &redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			DB:   cfg.Redis.DB,
		},
		Key:    cfg.Redis.Key,
		Logger: logger.WithField("FROM", "[REDIS]"),
	}

	return postgresConfig, kafkaConfig, redisConfig, nil
}

func parse() (*config, error) {
	k, err := getKafkaConfig()
	if err != nil {
		return nil, err
	}
	p, err := getPostgresConfig()
	if err != nil {
		return nil, err
	}
	r, err := getRedisConfig()
	if err != nil {
		return nil, err
	}
	return &config{
		Kafka:    k,
		Postgres: p,
		Redis:    r,
	}, nil
}
