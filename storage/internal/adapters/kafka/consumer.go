package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"storage/internal/app"
)

type Consumer struct {
	topics        []string
	handler       *Handler
	ctx           context.Context
	consumerGroup sarama.ConsumerGroup
	log           logrus.FieldLogger
}

func NewConsumer(app *app.App, cfg *Config) (*Consumer, error) {
	consumer := &Consumer{
		ctx: context.Background(),
	}
	handler := &Handler{
		app: app,
		log: cfg.Logger.WithField("FROM", "[KAFKA-HANDLER]"),
	}

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, InitConsumerConfig())
	if err != nil {
		cfg.Logger.Errorf("cannot create consumer group: %v", err)
		return nil, err
	}

	consumer.consumerGroup = consumerGroup
	consumer.topics = cfg.Topics
	consumer.handler = handler
	consumer.log = cfg.Logger
	return consumer, nil
}

func (c *Consumer) Run() error {
	for {
		if err := c.consumerGroup.Consume(c.ctx, c.topics, c.handler); err != nil {
			c.log.Errorf("cannot consume message: %v", err)
			return err
		}
		if c.ctx.Err() != nil {
			c.log.Errorf("got context error: %v", c.ctx.Err())
			return c.ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}
