package mq

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(writer *kafka.Writer) *Producer {
	return &Producer{
		writer: writer,
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
