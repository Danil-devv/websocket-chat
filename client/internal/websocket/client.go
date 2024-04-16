package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Client struct {
	conn     *websocket.Conn
	username string
}

func NewClient(host string, addr string, username string) *Client {
	u := url.URL{Scheme: "ws", Host: host, Path: addr}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	client := &Client{
		conn:     c,
		username: username,
	}
	return client
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) ReadMessage() (messageType int, msg Message, err error) {
	messageType, p, err := c.conn.ReadMessage()
	if err != nil {
		return messageType, Message{}, err
	}

	err = json.Unmarshal(p, &msg)
	if err != nil {
		return messageType, Message{}, err
	}

	return messageType, msg, err
}

func (c *Client) WriteMessage(messageType int, msg string) error {
	m := Message{
		Username: c.username,
		Text:     msg,
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(messageType, data)
}
