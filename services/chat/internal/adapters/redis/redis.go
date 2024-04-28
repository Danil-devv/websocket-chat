package redis

import (
	"chat/internal/domain"
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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

func (r *Repository) LoadMessages(ctx context.Context, count int) ([]domain.Message, error) {
	res := r.c.LRange(ctx, r.key, 0, max(0, int64(count-1)))
	data, err := res.Result()
	if err != nil {
		r.log.
			WithError(err).
			Error("cannot load messages")
		return nil, err
	}

	messages := make([]domain.Message, len(data))
	for i, v := range data {
		err = json.Unmarshal([]byte(v), &messages[len(messages)-i-1])
		if err != nil {
			r.log.
				WithError(err).
				WithField("message", v).
				Error("cannot unmarshall message")
			return nil, err
		}
	}
	return messages, nil
}
