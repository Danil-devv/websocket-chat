package config

import (
	"chat/internal/app"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type App struct {
	MessagesToLoad int
}

func getAppConfig() (*app.Config, error) {
	cfg, err := loadEnvAppConfig()
	if err != nil {
		return nil, err
	}
	return &app.Config{
		MessagesToLoad: cfg.MessagesToLoad,
	}, nil
}

func loadEnvAppConfig() (*App, error) {
	size, ok := os.LookupEnv("MESSAGES_TO_LOAD")
	if !ok {
		return nil, errors.New("cannot find 'MESSAGES_TO_LOAD' variable in environment")
	}
	messagesToLoad, err := strconv.Atoi(size)
	if err != nil {
		return nil, fmt.Errorf("%s: variable 'MESSAGES_TO_LOAD' must be integer", err.Error())
	}
	return &App{MessagesToLoad: messagesToLoad}, nil
}
