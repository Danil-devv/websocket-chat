package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/internal/adapters/websocket/syncmap"
	"server/internal/domain"
)

type App interface {
	SaveMessage(msg string, user string) error
	LoadLastMessages() ([]domain.Message, error)
}

type Server struct {
	srv http.Server
}

func NewServer(a App, cfg *Config, log logrus.FieldLogger) *Server {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  cfg.ReadBufferSize,
		WriteBufferSize: cfg.WriteBufferSize,
	}
	connections := syncmap.New()

	router := newRouter(a, upgrader, connections, log)

	return &Server{
		srv: http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: router,
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) GracefulShutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
