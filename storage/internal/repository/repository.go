package repository

import (
	"context"
	"log"
	"storage/internal/adapters/postgres"
	"storage/internal/adapters/redis"
	"storage/internal/domain"

	redisgh "github.com/redis/go-redis/v9"
)

type Repository struct {
	postgres *postgres.Repository
	redis    *redis.Repository
}

func New(pgConf *postgres.Config, redisConf *redisgh.Options) *Repository {
	return &Repository{
		postgres: postgres.NewRepository(pgConf),
		redis:    redis.NewRepository(redisConf),
	}
}

func (r *Repository) SaveMessage(ctx context.Context, message *domain.Message) error {
	err := r.postgres.SaveMessage(ctx, message)
	if err != nil {
		return err
	}
	go func() {
		err = r.redis.SaveMessage(ctx, message)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}
