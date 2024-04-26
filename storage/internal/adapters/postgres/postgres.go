package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"storage/internal/domain"
)

type Repository struct {
	pool *pgxpool.Pool
	log  logrus.FieldLogger
}

func NewRepository(conf *Config) *Repository {
	return &Repository{
		pool: conf.Pool,
		log:  conf.Logger,
	}
}

const saveMessageQuery = `INSERT INTO messages (username, data) VALUES ($1, $2);`

func (r *Repository) SaveMessage(ctx context.Context, message *domain.Message) error {
	r.log.Infof("got message: %+v", message)
	_, err := r.pool.Exec(ctx, saveMessageQuery, message.Username, message.Text)
	if err != nil {
		r.log.Errorf("cannot save message: %v", err)
		return err
	}
	r.log.Infof("successfully save message: %+v", message)
	return nil
}
