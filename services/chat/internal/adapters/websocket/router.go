package websocket

import (
	"chat/internal/adapters/websocket/syncmap"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

func newRouter(a App, u *websocket.Upgrader, connections *syncmap.ConnectionsMap, log logrus.FieldLogger) *http.ServeMux {
	r := &http.ServeMux{}
	r.HandleFunc("/api/v1/chat", createConnection(a, u, connections, log))
	return r
}
