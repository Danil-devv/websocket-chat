package websocket

type Message struct {
	Username string `json:"username" required:"true"`
	Text     string `json:"message" required:"true"`
}
