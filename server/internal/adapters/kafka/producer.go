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
	log   *logrus.Logger
	topic string
}

func NewProducer(cfg *ProducerConfig) (*Producer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.DefaultVersion
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for err = range producer.Errors() {
			cfg.Logger.WithError(err).Error("Failed to produce message")
		}
	}()

	return &Producer{
		conn:  producer,
		topic: cfg.Topic,
		log:   cfg.Logger,
	}, nil
}

func (p *Producer) SaveMessage(_ context.Context, message domain.Message) error {
	b, err := json.Marshal(message)
	p.log.Debugf("[KAFKA PRODUCER] trying to prouce message: %s", string(b))
	if err != nil {
		return err
	}
	p.conn.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   nil,
		Value: sarama.ByteEncoder(b),
	}
	return nil
}

func (p *Producer) Close() error {
	return p.conn.Close()
}
