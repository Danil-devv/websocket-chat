package main

import (
	"context"
	"errors"
	"fmt"
	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"server/internal/adapters/websocket"
	"server/internal/app"
	"server/internal/config"
	"server/internal/repository/postgres"
	"syscall"
)

var log *logrus.Logger

func main() {
	log = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.DebugLevel,
	}

	if err := godotenv.Load(); err != nil {
		log.WithError(err).Fatal("cannot load .env file")
	}

	pgxConfig, err := pgxpool.ParseConfig("postgres://postgres:postgres@db:5432/websocket-chat")

	if err != nil {
		log.WithError(err).Fatal("cannot parse config file")
	}

	pgxConfig.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(log),
		LogLevel: tracelog.LogLevelDebug,
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.WithError(err).Fatal("cannot create pgx pool with config")
	}
	defer pool.Close()

	repo := postgres.NewRepository(pool, log)

	appConfig, err := config.LoadApp()
	if err != nil {
		log.WithError(err).Fatal("cannot load .env file")
	}
	a := app.New(repo, appConfig)

	serverConfig, err := config.LoadServer()
	if err != nil {
		log.WithError(err).Fatal("cannot load .env file")
	}
	server := websocket.NewServer(a, serverConfig, log)

	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Infof("captured signal: %v", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		log.Infof("websocket server: start listening on :%s", serverConfig.Port)
		defer log.Infof("websocket server: close listening on :%s", serverConfig.Port)

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
		log.Infof("gracefully shutting down the server: %v", err)
	}

	log.Infof("start gracefull shutdown")
	err = server.GracefulShutdown(context.Background())
	if err != nil {
		log.Infof("cannot gracefully shutdown the server: %s", err.Error())
	} else {
		log.Infof("server was successfully shutted down")
	}
}
