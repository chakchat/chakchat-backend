package publish

import (
	"context"
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
	PublishForReceivers(ctx context.Context, users []uuid.UUID, typ string, data any) error
}

type UserEventPublisher struct {
	mq external.MqPublisher
}

func NewUserEventPublisher(mq external.MqPublisher) UserEventPublisher {
	return UserEventPublisher{
		mq: mq,
	}
}

func (p UserEventPublisher) PublishForReceivers(ctx context.Context, users []uuid.UUID, typ string, data any) error {
	if len(users) == 0 {
		return nil
	}

	e := UserEvent{
		Receivers: users,
		Type:      typ,
		Data:      data,
	}

	binE, err := json.Marshal(e)
	if err != nil {
		return err
	}

	return p.mq.Publish(ctx, binE)
}

type PublisherStub struct{}

func (PublisherStub) PublishForReceivers(ctx context.Context, users []uuid.UUID, typ string, data any) error {
	return nil
}
