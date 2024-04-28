package repository

import (
	"chat/internal/adapters/kafka"
	"chat/internal/adapters/postgres"
	"chat/internal/adapters/redis"
	"chat/internal/domain"
	"context"
	"time"
)

type Repository struct {
	postgres *postgres.Repository
	redis    *redis.Repository
	kafka    *kafka.Producer
}

func NewRepository(pgConf *postgres.Config, redisConf *redis.Config, kafkaConf *kafka.Config) (*Repository, error) {
	var (
		k   *kafka.Producer
		err error
	)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	retries := 15
	for retries > 0 {
		<-ticker.C
		k, err = kafka.NewProducer(kafkaConf)
		if err == nil {
			break
		}
		retries--
	}
	if err != nil {
		return nil, err
	}

	return &Repository{
		postgres: postgres.NewRepository(pgConf),
		redis:    redis.NewRepository(redisConf),
		kafka:    k,
	}, nil
}

func (r *Repository) SaveMessage(ctx context.Context, message domain.Message) error {
	return r.kafka.SaveMessage(ctx, message)
}

func (r *Repository) LoadMessages(ctx context.Context, count int) ([]domain.Message, error) {
	messages, err := r.redis.LoadMessages(ctx, count)
	if err == nil {
		return messages, nil
	}

	messages, err = r.postgres.LoadMessages(ctx, count)
	if err == nil {
		return messages, nil
	}
	return nil, err
}
