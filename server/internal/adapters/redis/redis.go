package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"server/internal/domain"
)

type Repository struct {
	c *redis.Client
}

func NewRepository(opt *redis.Options) *Repository {
	return &Repository{
		c: redis.NewClient(opt),
	}
}

func (r *Repository) LoadMessages(ctx context.Context, count int) ([]domain.Message, error) {
	res := r.c.LRange(ctx, "chat:messages", 0, max(0, int64(count-1)))
	data, err := res.Result()
	if err != nil {
		return nil, err
	}

	messages := make([]domain.Message, len(data))
	for i, v := range data {
		err = json.Unmarshal([]byte(v), &messages[len(messages)-i-1])
		if err != nil {
			return nil, err
		}
	}
	return messages, nil
}
