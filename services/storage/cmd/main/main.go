package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"storage/internal/adapters/kafka"
	"storage/internal/app"
	"storage/internal/config"
	"storage/internal/repository"
	"syscall"
)

const EnvFile = "example.env"

func main() {
	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			ForceColors:   true,
			PadLevelText:  true,
			FullTimestamp: true,
		},
		Hooks: logrus.LevelHooks{},
		Level: logrus.DebugLevel,
	}

	cfg, err := config.Get(logger, EnvFile)
	if err != nil {
		logger.Fatal(err)
	}

	repo := repository.New(cfg.Postgres, cfg.Redis, logger.WithField("FROM", "[REPOSITORY]"))
	a := app.NewApp(repo)
	consumer, err := kafka.NewConsumer(a, cfg.Kafka)
	if err != nil {
		logger.
			WithError(err).
			Fatal("cannot create kafka consumer")
	}

	// graceful shutdown
	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			logger.
				WithField("signal", s).
				Info("captured signal")
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		return consumer.Run()
	})

	if err = eg.Wait(); err != nil {
		logger.
			WithError(err).
			Info("gracefully shutting down the consumer")
	}

	if err = consumer.Close(); err != nil {
		logger.
			WithError(err).
			Error("failed to close consumer")
	}
}
