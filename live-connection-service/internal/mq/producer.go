package mq

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Brokers []string
	Topic   string
}
type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg ProducerConfig) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Brokers...),
			Topic:    cfg.Topic,
			Balancer: &kafka.Hash{},
			Async:    true,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) Send(ctx context.Context, msg any) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Value: jsonMsg,
	})

	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
