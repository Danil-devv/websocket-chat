package websocket

import (
	"chat/internal/adapters/websocket/syncmap"
	"chat/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

func createConnection(a App, u *websocket.Upgrader, c *syncmap.ConnectionsMap, log logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// --- OPEN NEW CONNECTION
		uid, conn, cancel, err := openNewConnection(log, u, w, r, c)
		defer cancel()
		if err != nil {
			return
		}
		// --- OPEN NEW CONNECTION

		// --- LOADING LAST MESSAGES
		log.WithField("uuid", uid.ID()).
			Info("start loading last messages")
		err = loadLastMessages(a, log, uid, conn)
		if err != nil {
			return
		}
		log.WithField("uuid", uid.ID()).
			Info("loading successfully finished")
		// --- LOADING LAST MESSAGES

		// --- LISTENING MESSAGES
		log.WithField("uuid", uid.ID()).
			Info("start listening messages")
		for {
			messageType, data, err := conn.ReadMessage()
			log.WithField("uuid", uid.ID()).
				WithField("data", string(data)).
				Info("got message")

			if err != nil || messageType == websocket.CloseMessage {
				log.WithError(err).
					WithField("uuid", uid.ID()).
					WithField("message type", messageType).
					Info("closing the connection")
				return
			}

			err = validateMessage(data)
			if err != nil {
				log.WithError(err).
					WithField("uuid", uid.ID()).
					Info("message doesnt pass validation, receiving an error message")

				msg := domain.Message{
					Username: "WRONG MESSAGE ERROR",
					Text:     err.Error(),
				}

				data, err := json.Marshal(msg)
				if err != nil {
					log.WithError(err).
						WithField("uuid", uid.ID()).
						WithField("msg", msg).
						Error("cannot marshal data to json")
					continue
				}

				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.WithError(err).
						WithField("uuid", uid.ID()).
						Error("cannot send message to client")
					return
				}
				continue
			}

			go saveAndSendMessage(data, c, a, log)
		}
		// --- LISTENING MESSAGES
	}
}

func openNewConnection(
	log logrus.FieldLogger, u *websocket.Upgrader,
	w http.ResponseWriter, r *http.Request,
	c *syncmap.ConnectionsMap,
) (uid uuid.UUID, conn *websocket.Conn, cancelFunc func(), err error) {
	uid = uuid.New()
	log.WithField("uuid", uid.ID()).
		Info("trying to open new websocket connection")
	conn, err = u.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).
			WithField("uuid", uid.ID()).
			Error("cannot upgrade connection to websocket")
		return
	}
	c.Store(conn)
	log.WithField("uuid", uid.ID()).
		Info("store the connection")
	return uid, conn, func() {
		c.Delete(conn)
		log.WithField("uuid", uid.ID()).
			Info("delete the connection")
		err := conn.Close()
		if err != nil {
			log.WithError(err).
				WithField("uuid", uid.ID()).
				Error("cannot close connection")
		}
	}, nil
}

func loadLastMessages(a App, log logrus.FieldLogger, uid uuid.UUID, conn *websocket.Conn) error {
	messages, err := a.LoadLastMessages()
	if err != nil {
		defer func() {
			err = conn.WriteMessage(
				websocket.CloseInternalServerErr,
				[]byte(fmt.Sprintf("cannot load last messages: %s", err.Error())),
			)
			if err != nil {
				log.WithError(err).
					WithField("uuid", uid.ID()).
					Error("cannot send message to client")
			}
		}()

		log.WithError(err).
			WithField("uuid", uid.ID()).
			Error("cannot load last messages")
		return err
	}

	log.WithField("uuid", uid.ID()).
		Info("start sending last messages")
	err = saveLastMessages(messages, log, uid, conn)
	log.WithField("uuid", uid.ID()).
		Info("last messages have been sent")
	return err
}

func saveLastMessages(messages []domain.Message, log logrus.FieldLogger, uid uuid.UUID, conn *websocket.Conn) error {
	for _, m := range messages {
		data, err := json.Marshal(m)
		if err != nil {
			log.WithError(err).
				WithField("uuid", uid.ID()).
				Info("cannot marshal data to json")
			continue
		}
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.WithError(err).
				WithField("uuid", uid.ID()).
				Error("cannot send message to client")
			return err
		}
	}
	return nil
}

func saveAndSendMessage(data []byte, c *syncmap.ConnectionsMap, a App, l logrus.FieldLogger) {
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
