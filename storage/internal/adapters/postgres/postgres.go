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
	r.log.
		WithField("message", message).
		Info("got message")
	_, err := r.pool.Exec(ctx, saveMessageQuery, message.Username, message.Text)
	if err != nil {
		r.log.
			WithError(err).
			Error("cannot save message")
		return err
	}
	r.log.
		WithField("message", message).
		Info("successfully save message")
	return nil
}
