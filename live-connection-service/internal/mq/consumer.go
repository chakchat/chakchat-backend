package mq

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type ConsumerConf struct {
	Brokers []string
	Topic   string
	GroupID string
}

type Consumer struct {
	reader   *kafka.Reader
	handler  func(ctx context.Context, msg kafka.Message) error
	shutdown chan struct{}
}

func NewConsumer(reader *ConsumerConf) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Topic:          reader.Topic,
			Brokers:        reader.Brokers,
			GroupID:        reader.GroupID,
			StartOffset:    kafka.LastOffset,
			CommitInterval: 0,
		}),
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
				msg, err := c.reader.FetchMessage(ctx)
				if err != nil {
					if err == context.Canceled {
						return
					}
					continue
				}
				processErr := c.handler(ctx, msg)
				if processErr != nil {
					continue
				}

				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					continue
				}
			}
		}
	}()
}

func (c *Consumer) Stop() {
	close(c.shutdown)
	_ = c.reader.Close()
}
