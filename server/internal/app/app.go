package app

import (
	"context"
	"server/internal/config"
	"server/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name=Repository
type Repository interface {
	SaveMessage(ctx context.Context, message domain.Message) error
	LoadMessages(ctx context.Context, count int) ([]domain.Message, error)
}

type App struct {
	messagesToLoad int
	repo           Repository
}

func New(r Repository, conf *config.App) *App {
	return &App{
		repo:           r,
		messagesToLoad: conf.MessagesToLoad,
	}
}

func (a *App) SaveMessage(msg string, user string) error {
	err := a.repo.SaveMessage(
		context.Background(),
		domain.Message{Username: user, Text: msg},
	)

	if err != nil {
		return newAppError(err)
	}
	return nil
}

func (a *App) LoadLastMessages() ([]domain.Message, error) {
	messages, err := a.repo.LoadMessages(
		context.Background(),
		a.messagesToLoad,
	)

	if err != nil {
		return nil, newAppError(err)
	}
	return messages, nil
}
