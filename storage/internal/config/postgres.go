package config

import (
	"fmt"
	"os"
)

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func getPostgresConfig() (Postgres, error) {
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
