package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"server/internal/adapters/websocket"
	"server/internal/app"
	"server/internal/config"
	"server/internal/repository"
	"syscall"
)

func main() {
	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			ForceColors:   true,
			PadLevelText:  true,
			FullTimestamp: true,
		},
		Hooks: logrus.LevelHooks{},
		Level: logrus.DebugLevel,
	}

	cfg, err := config.Get(logger, "example.env")
	if err != nil {
		logger.WithError(err).Fatal("cannot parse config")
	}

	repo, err := repository.NewRepository(cfg.Postgres, cfg.Redis, cfg.Kafka)
	if err != nil {
		logger.WithError(err).Fatal("cannot create repository")
	}
	a := app.New(repo, cfg.App)
	server := websocket.NewServer(a, cfg.Server, logger)

	// graceful shutdown
	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			logger.WithField("signal", s).Info("captured signal")
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		logger.WithField("port", cfg.Server.Port).Info("websocket server: start listening")
		defer logger.WithField("port", cfg.Server.Port).Infof("websocket server: close listening")

		errCh := make(chan error)

		go func() {
			err := server.ListenAndServe()
			if !errors.Is(err, http.ErrServerClosed) && err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("websocket server can't listen and serve requests: %s", err.Error())
		}
	})

	if err = eg.Wait(); err != nil {
		logger.WithError(err).Info("gracefully shutting down the server")
	}

	logger.Info("start gracefull shutdown")
	err = server.GracefulShutdown(context.Background())
	if err != nil {
		logger.Infof("cannot gracefully shutdown the server: %s", err.Error())
	} else {
		logger.Infof("server was successfully shutted down")
	}
}
