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

type Publisher interface {
	PublishForUsers(users []uuid.UUID, ev Event)
}

type UserEventPublisher struct {
	mq external.MqPublisher
}

func NewUserEventPublisher(mq external.MqPublisher) UserEventPublisher {
	return UserEventPublisher{
		mq: mq,
	}
}

func (p UserEventPublisher) PublishForUsers(users []uuid.UUID, ev Event) {
	if len(users) == 0 {
		return
	}

	userEvent := UserEvent{
		Users: users,
		Data:  ev,
	}

	p.mq.Publish(userEvent)
}
