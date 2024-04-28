package syncmap

import (
	"github.com/gorilla/websocket"
	"sync"
)

type ConnectionsMap struct {
	mx *sync.RWMutex
	m  map[*websocket.Conn]struct{}
}

func New() *ConnectionsMap {
	return &ConnectionsMap{
		mx: &sync.RWMutex{},
		m:  make(map[*websocket.Conn]struct{}),
	}
}

func (c *ConnectionsMap) LoadAllConnections() <-chan *websocket.Conn {
	c.mx.Lock()

	ch := make(chan *websocket.Conn, len(c.m))
	go func() {
		defer func() {
			c.mx.Unlock()
			close(ch)
		}()
		for conn := range c.m {
			ch <- conn
		}
	}()

	return ch
}

func (c *ConnectionsMap) Store(key *websocket.Conn) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	c.m[key] = struct{}{}
}

func (c *ConnectionsMap) Delete(key *websocket.Conn) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	delete(c.m, key)
}
