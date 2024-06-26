package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"storage/internal/app"
	"storage/internal/domain"
)

type Handler struct {
	app *app.App
	log logrus.FieldLogger
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				h.log.Infoln("message channel was closed")
				return nil
			}

			h.log.
				WithField("message", string(message.Value)).
				Info("message claimed")
			msg := &domain.Message{}
			err := json.Unmarshal(message.Value, msg)
			if err != nil {
				h.log.
					WithError(err).
					WithField("message.value", message.Value).
					Errorf("cannot unmarshal message")
				return err
			}

			h.log.Infoln("saving message")
			err = h.app.SaveMessage(session.Context(), msg)
			if err != nil {
				h.log.
					WithError(err).
					WithField("message", msg).
					Error("cannot save message")
				return nil
			}

			h.log.
				WithField("message", message).
				Info("message successfully saved, marking message")
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			h.log.Infoln("session is done")
			session.Commit()
			return nil
		}
	}
}
