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
	r.log.Infof("got message: %+v", message)
	data, err := json.Marshal(message)
	if err != nil {
		r.log.WithField("message", message).
			Errorf("failed to marshal message: %v", err)
		return err
	}
	r.log.Infof("saving message: %s", data)
	r.c.LPush(ctx, r.key, string(data))
	return nil
}
