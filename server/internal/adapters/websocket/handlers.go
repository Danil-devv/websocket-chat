package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/internal/adapters/websocket/syncmap"
	"server/internal/domain"
)

func createConnection(a App, u *websocket.Upgrader, c *syncmap.ConnectionsMap, l *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := u.Upgrade(w, r, nil)

		uid := uuid.New()
		l.WithField("uuid", uid.ID()).Debug("trying to open new websocket connection")
		defer func() {
			err := conn.Close()
			if err != nil {
				l.WithField("uuid", uid.ID()).WithError(err).Error("cannot close connection")
			}
		}()
		if err != nil {
			l.WithField("uuid", uid.ID()).WithError(err).Error("cannot upgrade connection to websocket")
			return
		}

		c.Store(conn)
		l.WithField("uuid", uid.ID()).Debug("store the connection")
		defer func() {
			c.Delete(conn)
			l.WithField("uuid", uid.ID()).Debug("delete the connection")
		}()

		l.WithField("uuid", uid.ID()).Debug("start loading last messages")
		messages, err := a.LoadLastMessages()
		if err != nil {
			l.WithField("uuid", uid.ID()).WithError(err).Error("internal error: cannot load last messages")
			err := conn.WriteMessage(
				websocket.CloseInternalServerErr,
				[]byte(fmt.Sprintf("cannot load last messages: %s", err.Error())),
			)
			if err != nil {
				l.WithField("uuid", uid.ID()).WithError(err).Error("cannot send message to client")
			}
			return
		}
		l.WithField("uuid", uid.ID()).Debug("loading successfully finished")

		l.WithField("uuid", uid.ID()).Debug("start sending last messages")
		for _, m := range messages {
			data, err := json.Marshal(m)
			if err != nil {
				l.WithField("uuid", uid.ID()).WithError(err).Info("cannot marshal data to json")
				continue
			}
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				l.WithField("uuid", uid.ID()).WithError(err).Error("cannot send message to client")
				return
			}
		}
		l.WithField("uuid", uid.ID()).Debug("last messages have been sent")

		l.WithField("uuid", uid.ID()).Debug("start listening messages")
		for {
			messageType, data, err := conn.ReadMessage()
			l.WithField("uuid", uid.ID()).WithField("data", string(data)).Debug("got message")

			if err != nil || messageType == websocket.CloseMessage {
				l.WithField("uuid", uid.ID()).
					WithError(err).
					WithField("message type", messageType).
					Info("closing the connection")
				return
			}

			err = validateMessage(data)
			if err != nil {
				l.WithField("uuid", uid.ID()).
					Debug("message doesnt pass validation, receiving an error message")

				msg := domain.Message{
					Username: "WRONG MESSAGE ERROR",
					Text:     err.Error(),
				}

				data, err := json.Marshal(msg)
				if err != nil {
					l.WithField("uuid", uid.ID()).WithError(err).Error("cannot marshal data to json")
					continue
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					l.WithField("uuid", uid.ID()).WithError(err).Error("cannot send message to client")
					return
				}
				continue
			}

			go saveAndSendMessage(data, c, a, l)
		}
	}
}

func saveAndSendMessage(data []byte, c *syncmap.ConnectionsMap, a App, l *logrus.Logger) {
	msg := domain.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		l.WithError(err).WithField("data", string(data)).Error("cannot unmarshal data")
		return
	}

	err = a.SaveMessage(msg.Text, msg.Username)
	if err != nil {
		l.WithError(err).WithField("message", msg).Error("cannot save message")
		return
	}

	ch := c.LoadAllConnections()
	for conn := range ch {
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			l.WithError(err).WithField("data", string(data)).Error("cannot send the message")
		}
	}
}

func validateMessage(data []byte) error {
	msg := domain.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	if msg.Text == "" {
		return errors.New("message text must be non-empty")
	}

	if len(msg.Username) < 3 {
		return errors.New("username length must be at least 3 characters")
	}

	return nil
}
