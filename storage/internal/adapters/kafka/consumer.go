package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"storage/internal/app"
)

type Consumer struct {
	cfg           *Config
	handler       *Handler
	ctx           context.Context
	consumerGroup sarama.ConsumerGroup
}

func NewConsumer(ctx context.Context, app app.App, cfg *Config) (*Consumer, error) {
	consumer := &Consumer{
		ctx: ctx,
	}
	handler := &Handler{app: app}

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, InitConsumerConfig())
	if err != nil {
		log.Printf("Error creating consumerGroup %v", err)
		return nil, err
	}

	consumer.consumerGroup = consumerGroup
	consumer.cfg = cfg
	consumer.handler = handler
	return consumer, nil
}

func (c *Consumer) Run() error {
	for {
		if err := c.consumerGroup.Consume(c.ctx, c.cfg.Topics, c.handler); err != nil {
			log.Fatalf("Error from consumer: %v", err)
			return err
		}
		if c.ctx.Err() != nil {
			return c.ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}
