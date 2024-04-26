package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"server/internal/domain"
)

type Repository struct {
	r *redis.Client
}

func NewRepository(opt *redis.Options) *Repository {
	return &Repository{
		r: redis.NewClient(opt),
	}
}

func (c *Repository) LoadMessages(ctx context.Context, count int) ([]domain.Message, error) {
	res := c.r.LRange(ctx, "chat:messages", 0, int64(count))
	data, err := res.Result()
	if err != nil {
		return nil, err
	}

	messages := make([]domain.Message, count)
	for i, v := range data {
		err = json.Unmarshal([]byte(v), &messages[i])
		if err != nil {
			return nil, err
		}
	}
	return messages, nil
}
