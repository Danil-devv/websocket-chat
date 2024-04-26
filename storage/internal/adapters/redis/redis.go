package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"storage/internal/domain"
)

type Repository struct {
	c *redis.Client
}

func NewRepository(opt *redis.Options) *Repository {
	return &Repository{
		c: redis.NewClient(opt),
	}
}

func (r *Repository) SaveMessage(ctx context.Context, message *domain.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	r.c.LPush(ctx, "chat:messages", string(data))
	return nil
}
