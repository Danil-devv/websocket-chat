package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"storage/internal/domain"
)

type Repository struct {
	c   *redis.Client
	key string
	log logrus.FieldLogger
}

func NewRepository(cfg *Config) *Repository {
	return &Repository{
		c:   redis.NewClient(cfg.Opt),
		key: cfg.Key,
		log: cfg.Logger,
	}
}

func (r *Repository) SaveMessage(ctx context.Context, message *domain.Message) error {
	r.log.
		WithField("message", message).
		Info("got message")
	data, err := json.Marshal(message)
	if err != nil {
		r.log.
			WithError(err).
			WithField("message", message).
			Error("failed to marshal message")
		return err
	}
	r.log.
		WithField("message", message).
		Info("saving message")
	r.c.LPush(ctx, r.key, string(data))
	return nil
}
