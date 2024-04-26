package kafka

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log"
	"storage/internal/app"
	"storage/internal/domain"
)

type Handler struct {
	app *app.App
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
				log.Println("message channel was closed")
				return nil
			}

			msg := &domain.Message{}
			err := json.Unmarshal(message.Value, msg)
			if err != nil {
				log.Println("error unmarshalling message: ", err.Error())
				return err
			}

			err = h.app.SaveMessage(session.Context(), msg)
			if err != nil {
				log.Printf("error {%s} while handling message message: %+v", err, message)
				return nil
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			session.Commit()
			return nil
		}
	}
}
