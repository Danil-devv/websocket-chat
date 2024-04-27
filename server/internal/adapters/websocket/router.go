package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/internal/adapters/websocket/syncmap"
)

func newRouter(a App, u *websocket.Upgrader, connections *syncmap.ConnectionsMap, log logrus.FieldLogger) *http.ServeMux {
	r := &http.ServeMux{}
	r.HandleFunc("/api/v1/chat", createConnection(a, u, connections, log))
	return r
}
