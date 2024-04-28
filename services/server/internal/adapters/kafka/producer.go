package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"server/internal/domain"
	"time"
)

type Producer struct {
	conn  sarama.AsyncProducer
	log   logrus.FieldLogger
	topic string
}

func NewProducer(cfg *Config) (*Producer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, config)
	if err != nil {
		cfg.Logger.WithError(err).Error("cannot create new kafka producer")
		return nil, err
	}

	go func() {
		for err = range producer.Errors() {
			cfg.Logger.WithError(err).Error("failed to produce message")
		}
	}()

	return &Producer{
		conn:  producer,
		topic: cfg.Topic,
		log:   cfg.Logger,
	}, nil
}

func (p *Producer) SaveMessage(_ context.Context, message domain.Message) error {
	p.log.
		WithField("message", message).
		Info("trying to produce message")
	b, err := json.Marshal(message)
	if err != nil {
		p.log.WithError(err).Error("cannot marshal message")
		return err
	}
	p.conn.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   nil,
		Value: sarama.ByteEncoder(b),
	}
	p.log.
		WithField("message", message).
		Info("message was produced")
	return nil
}

func (p *Producer) Close() error {
	return p.conn.Close()
}
