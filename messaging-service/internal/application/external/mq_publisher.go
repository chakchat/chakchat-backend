package external

import "context"

type Event any

type MqPublisher interface {
	// It doesn't return a value to guarantee that event will be published.
	// Eventual consistency is possible.
	Publish(context.Context, []byte) error
}
