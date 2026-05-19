package mq

import (
	"context"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader   *kafka.Reader
	handler  func(ctx context.Context, msg kafka.Message) error
	shutdown chan struct{}
}

func NewConsumer(reader *kafka.Reader) *Consumer {
	return &Consumer{
		reader:   reader,
		shutdown: make(chan struct{}),
	}
}

func (c *Consumer) Start(ctx context.Context, handler func(ctx context.Context, msg kafka.Message) error) {
	c.handler = handler
	go func() {
		for {
			select {
			case <-c.shutdown:
				return
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						log.Printf("kafka message reading: %v", err)
						return
					}
					log.Printf("Can't read message from kafka: %v", err)
					continue
				}
				processErr := c.handler(ctx, msg)
				if processErr != nil {
					log.Printf("Error to handle kafka message: %s", err)
					continue
				}
				log.Printf("Successfully handle kafka message. Ready to commit")
			}
		}
	}()
}

func (c *Consumer) Stop() {
	close(c.shutdown)
	_ = c.reader.Close()
}
