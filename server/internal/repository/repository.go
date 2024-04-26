package repository

import (
	"context"
	redisgh "github.com/redis/go-redis/v9"
	"server/internal/adapters/kafka"
	"server/internal/adapters/postgres"
	"server/internal/adapters/redis"
	"server/internal/domain"
)

type Repository struct {
	postgres *postgres.Repository
	redis    *redis.Repository
	kafka    *kafka.Producer
}

func NewRepository(pgConf *postgres.Config, redisConf *redisgh.Options, kafkaConf *kafka.ProducerConfig) (*Repository, error) {
	k, err := kafka.NewProducer(kafkaConf)
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
