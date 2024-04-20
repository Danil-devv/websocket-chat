package syncmap

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectionsMap_Store(t *testing.T) {
	type testcase struct {
		conns []*websocket.Conn
	}

	tests := []testcase{
		{
			conns: make([]*websocket.Conn, 10),
		}, {
			conns: make([]*websocket.Conn, 1),
		}, {
			conns: make([]*websocket.Conn, 0),
		}, {
			conns: make([]*websocket.Conn, 100),
		},
	}

	for i := 0; i < len(tests); i++ {
		for j := 0; j < len(tests[i].conns); j++ {
			tests[i].conns[j] = new(websocket.Conn)
		}
	}

	for _, test := range tests {
		connMap := New()
		exists := make(map[*websocket.Conn]bool)
		for _, conn := range test.conns {
			connMap.Store(conn)
		}

		ch := connMap.LoadAllConnections()
		for conn := range ch {
			if exists[conn] {
				t.Errorf("duplicate of connections: %p", conn)
			}
			exists[conn] = true
		}
		assert.Equal(t, len(test.conns), len(exists))
	}
}

func TestConnectionsMap_Delete(t *testing.T) {
	type testcase struct {
		conns []*websocket.Conn
	}

	tests := []testcase{
		{
			conns: make([]*websocket.Conn, 10),
		}, {
			conns: make([]*websocket.Conn, 1),
		}, {
			conns: make([]*websocket.Conn, 0),
		}, {
			conns: make([]*websocket.Conn, 100),
		},
	}

	for i := 0; i < len(tests); i++ {
		for j := 0; j < len(tests[i].conns); j++ {
			tests[i].conns[j] = new(websocket.Conn)
		}
	}

	for _, test := range tests {
		connMap := New()
		for _, conn := range test.conns {
			connMap.Store(conn)
		}

		ch := connMap.LoadAllConnections()
		for conn := range ch {
			connMap.Delete(conn)
		}

		ch = connMap.LoadAllConnections()
		count := 0
		for range ch {
			count++
		}

		assert.Equal(t, 0, count)
	}
}
