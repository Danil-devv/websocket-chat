package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/internal/websocket/syncmap"
)

func newRouter(a App, u *websocket.Upgrader, connections *syncmap.ConnectionsMap, l *logrus.Logger) *http.ServeMux {
	r := &http.ServeMux{}
	r.HandleFunc("/api/v1/chat", createConnection(a, u, connections, l))
	return r
}
