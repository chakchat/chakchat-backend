package messages

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
	handler  func(ctx context.Context, msg kafka.Message)
	shutdown chan struct{}
}

func NewConsumer(reader *ConsumerConf) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Topic:       reader.Topic,
			Brokers:     reader.Brokers,
			GroupID:     reader.GroupID,
			StartOffset: kafka.LastOffset,
		}),
		shutdown: make(chan struct{}),
	}
}

func (c *Consumer) Start(ctx context.Context, handler func(ctx context.Context, msg kafka.Message)) {
	c.handler = handler
	go func() {
		for {
			select {
			case <-c.shutdown:
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					continue
				}
				go c.handler(ctx, msg)
			}
		}
	}()
}

func (c *Consumer) Stop() {
	close(c.shutdown)
	_ = c.reader.Close()
}
