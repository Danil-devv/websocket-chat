package config

import (
	"fmt"
	"os"
	"strconv"
)

type Redis struct {
	Host string
	Port string
	Key  string
	DB   int
}

func getRedisConfig() (Redis, error) {
	host, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		return Redis{}, fmt.Errorf("REDIS_HOST environment variable not set")
	}
	port, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		return Redis{}, fmt.Errorf("REDIS_PORT environment variable not set")
	}
	dbString, ok := os.LookupEnv("REDIS_DB")
	if !ok {
		return Redis{}, fmt.Errorf("REDIS_DB environment variable not set")
	}
	db, err := strconv.Atoi(dbString)
	if err != nil {
		return Redis{}, fmt.Errorf("REDIS_DB environment variable not integer")
	}
	key, ok := os.LookupEnv("REDIS_KEY")
	if !ok {
		return Redis{}, fmt.Errorf("REDIS_KEY environment variable not set")
	}
	return Redis{
		Host: host,
		Port: port,
		Key:  key,
		DB:   db,
	}, nil
}
