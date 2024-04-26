package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"storage/internal/adapters/postgres"
	"storage/internal/adapters/redis"
	"storage/internal/domain"
)

type Repository struct {
	postgres *postgres.Repository
	redis    *redis.Repository
	log      logrus.FieldLogger
}

func New(pgConf *postgres.Config, redisConf *redis.Config, log logrus.FieldLogger) *Repository {
	return &Repository{
		postgres: postgres.NewRepository(pgConf),
		redis:    redis.NewRepository(redisConf),
		log:      log,
	}
}

func (r *Repository) SaveMessage(ctx context.Context, message *domain.Message) error {
	err := r.postgres.SaveMessage(ctx, message)
	if err != nil {
		r.log.Errorf("cannot save message to postgres: %v", err)
		return err
	}
	go func() {
		err = r.redis.SaveMessage(ctx, message)
		if err != nil {
			r.log.Errorf("cannot save message to redis: %v", err)
		}
	}()
	return nil
}
