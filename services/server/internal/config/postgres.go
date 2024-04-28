package config

import (
	"context"
	"fmt"
	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
	"os"
	"server/internal/adapters/postgres"
)

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func getPostgresConfig(logger logrus.FieldLogger) (*postgres.Config, error) {
	config, err := loadEnvPostgresConfig()
	if err != nil {
		return nil, err
	}

	pgxConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		),
	)
	if err != nil {
		logger.
			WithError(err).
			Error("cannot parse postgres config")
		return nil, err
	}
	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(logger.WithField("FROM", "[PGX-POOL]")),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		logger.
			WithError(err).
			Error("cannot create pgx pool")
		return nil, err
	}
	postgresConfig := &postgres.Config{
		Pool:   pool,
		Logger: logger.WithField("FROM", "[POSTGRES]"),
	}
	return postgresConfig, nil
}

func loadEnvPostgresConfig() (Postgres, error) {
	username, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		return Postgres{}, fmt.Errorf("POSTGRES_USER environment variable not set")
	}
	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return Postgres{}, fmt.Errorf("POSTGRES_PASSWORD environment variable not set")
	}
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return Postgres{}, fmt.Errorf("POSTGRES_HOST environment variable not set")
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		return Postgres{}, fmt.Errorf("POSTGRES_PORT environment variable not set")
	}
	database, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		return Postgres{}, fmt.Errorf("POSTGRES_DB environment variable not set")
	}
	return Postgres{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}, nil
}
