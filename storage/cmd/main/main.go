package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"storage/internal/adapters/kafka"
	"storage/internal/app"
	"strings"
	"syscall"
)

func main() {
	a := app.NewApp(_)

	kafkaConfig := &kafka.Config{
		Brokers: strings.Split("kafka1:29092,kafka2:29093,kafka3:29094", ","),
		Topics:  []string{"ts.2s.2"},
		GroupID: "1",
	}
	consumer, err := kafka.NewConsumer(a, kafkaConfig)

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
