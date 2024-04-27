package main

import (
	"context"
	"fmt"
	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"storage/internal/adapters/kafka"
	"storage/internal/adapters/postgres"
	rds "storage/internal/adapters/redis"
	"storage/internal/app"
	"storage/internal/config"
	"storage/internal/repository"
	"strings"
	"syscall"
)

func main() {
	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			ForceColors:  true,
			PadLevelText: true,
		},
		Hooks: logrus.LevelHooks{},
		Level: logrus.DebugLevel,
	}

	postgresConfig, kafkaConfig, redisConfig, err := getConfigs(logger)
	if err != nil {
		logger.Fatal(err)
	}

	repo := repository.New(postgresConfig, redisConfig, logger.WithField("FROM", "[REPOSITORY]"))
	a := app.NewApp(repo)
	consumer, err := kafka.NewConsumer(a, kafkaConfig)
	if err != nil {
		logger.Fatalf("cannot create kafka consumer: %v", err)
	}

	// graceful shutdown
	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			logger.Printf("captured signal: %v", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		return consumer.Run()
	})

	if err = eg.Wait(); err != nil {
		logger.Printf("gracefully shutting down the consumer: %v", err)
	}

	if err = consumer.Close(); err != nil {
		logger.Printf("failed to close consumer: %v", err)
	}
}

func getConfigs(logger *logrus.Logger) (*postgres.Config, *kafka.Config, *rds.Config, error) {
	cfg, err := config.Get()
	if err != nil {
		logger.Errorf("cannot get config: %v", err)
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
		logger.Errorf("cannot parse postgres config: %v", err)
		return nil, nil, nil, err
	}
	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(logger.WithField("FROM", "[PGX-POOL]")),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		logger.Errorf("cannot create pool: %v", err)
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
