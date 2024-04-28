package app

import (
	"context"
	"storage/internal/domain"
)

type MessageSaver interface {
	SaveMessage(ctx context.Context, msg *domain.Message) error
}

type App struct {
	repository MessageSaver
}

func NewApp(repository MessageSaver) *App {
	return &App{repository: repository}
}

func (a *App) SaveMessage(ctx context.Context, msg *domain.Message) error {
	return a.repository.SaveMessage(ctx, msg)
}
