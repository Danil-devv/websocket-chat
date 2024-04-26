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
	"log"
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
	// logger settings
	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			ForceColors:  true,
			PadLevelText: true,
		},
		Hooks: logrus.LevelHooks{},
		Level: logrus.DebugLevel,
	}

	postgresConfig, kafkaConfig, redisConfig, err := getConfigs()
	if err != nil {
		log.Fatal(err)
	}

	kafkaConfig.Logger = logger.WithField("FROM", "[KAFKA-CONSUMER]")
	postgresConfig.Logger = logger.WithField("FROM", "[POSTGRES]")
	redisConfig.Logger = logger.WithField("FROM", "[REDIS]")

	repo := repository.New(postgresConfig, redisConfig, logger.WithField("FROM", "[REPOSITORY]"))
	a := app.NewApp(repo)
	consumer, err := kafka.NewConsumer(a, kafkaConfig)

	// graceful shutdown
	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		return consumer.Run()
	})

	if err = eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the consumer: %v", err)
	}

	if err = consumer.Close(); err != nil {
		log.Printf("failed to close consumer: %v", err)
	}
}

func getConfigs() (*postgres.Config, *kafka.Config, *rds.Config, error) {
	cfg, err := config.Get()
	if err != nil {
		log.Println(err)
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
		log.Println("cannot parse config file")
		return nil, nil, nil, err
	}
	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(log.Logger{}),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Println("cannot create pgx pool with config")
		return nil, nil, nil, err
	}
	postgresConfig := &postgres.Config{
		Pool: pool,
	}

	kafkaConfig := &kafka.Config{
		Brokers: strings.Split(cfg.Kafka.Brokers, ","),
		Topics:  strings.Split(cfg.Kafka.Topics, ","),
		GroupID: cfg.Kafka.GroupID,
	}
	redisConfig := &rds.Config{
		Opt: &redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			DB:   cfg.Redis.DB,
		},
		Key: cfg.Redis.Key,
	}

	return postgresConfig, kafkaConfig, redisConfig, nil
}
