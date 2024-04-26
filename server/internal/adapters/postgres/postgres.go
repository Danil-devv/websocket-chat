package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"server/internal/domain"
)

type Repository struct {
	pool   *pgxpool.Pool
	logger logrus.FieldLogger
}

func NewRepository(conf *Config) *Repository {
	return &Repository{
		pool:   conf.Pool,
		logger: conf.Logger,
	}
}

const saveMessageQuery = `INSERT INTO messages (username, data) VALUES ($1, $2);`

func (r *Repository) SaveMessage(ctx context.Context, message domain.Message) error {
	_, err := r.pool.Exec(ctx, saveMessageQuery, message.Username, message.Text)
	if err != nil {
		return newPostgresError(err)
	}
	return nil
}

const loadMessagesQuery = `SELECT username, data FROM
    (SELECT * FROM
        messages
        ORDER BY id DESC LIMIT $1)
ORDER BY id;`

func (r *Repository) LoadMessages(ctx context.Context, count int) ([]domain.Message, error) {
	rows, err := r.pool.Query(ctx, loadMessagesQuery, count)
	if err != nil {
		return nil, newPostgresError(err)
	}
	defer rows.Close()

	res := make([]domain.Message, 0)
	for rows.Next() {
		msg := domain.Message{}
		err = rows.Scan(&msg.Username, &msg.Text)
		if err != nil {
			return nil, newPostgresError(err)
		}
		res = append(res, msg)
	}
	return res, nil
}
