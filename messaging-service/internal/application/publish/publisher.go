package publish

import "github.com/google/uuid"

type Event any

type Publisher interface {
	PublishForUsers([]uuid.UUID, Event) error
}
