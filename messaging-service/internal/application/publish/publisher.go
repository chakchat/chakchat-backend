package publish

import (
	"encoding/json"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/google/uuid"
)

type UserEvent struct {
	Receivers []uuid.UUID `json:"receivers"`
	Type      string      `json:"type"`
	Data      any         `json:"data"`
}

type Publisher interface {
	PublishForReceivers(users []uuid.UUID, typ string, data any)
}

type UserEventPublisher struct {
	mq external.MqPublisher
}

func NewUserEventPublisher(mq external.MqPublisher) UserEventPublisher {
	return UserEventPublisher{
		mq: mq,
	}
}

func (p UserEventPublisher) PublishForReceivers(users []uuid.UUID, typ string, data any) {
	if len(users) == 0 {
		return
	}

	e := UserEvent{
		Receivers: users,
		Type:      typ,
		Data:      data,
	}

	binE, err := json.Marshal(e)
	if err != nil {
		panic(err) // TODO: use smth like outbox and here return an error
	}

	p.mq.Publish(binE)
}

type PublisherStub struct{}

func (PublisherStub) PublishForReceivers(users []uuid.UUID, typ string, data any) {}
