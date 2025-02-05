package publish

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/google/uuid"
)

type Event any

type UserEvent struct {
	Users []uuid.UUID `json:"users"`
	Data  Event       `json:"data"`
}

type Publisher struct {
	mq external.MqPublisher
}

func NewPublisher(mq external.MqPublisher) Publisher {
	return Publisher{
		mq: mq,
	}
}

func (p Publisher) PublishForUsers(users []uuid.UUID, ev Event) error {
	userEvent := UserEvent{
		Users: users,
		Data:  ev,
	}
	return p.mq.Publish(userEvent)
}
