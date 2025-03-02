package configuration

import "github.com/chakchat/chakchat-backend/messaging-service/internal/rest/handlers/chat"

type Handlers struct {
	PersonalChat *chat.PersonalChatHandler
}

func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		PersonalChat: chat.NewPersonalChatHandler(services.PersonalChat),
	}
}
