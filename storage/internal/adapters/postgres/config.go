package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Pool   *pgxpool.Pool
	Logger logrus.FieldLogger
}
