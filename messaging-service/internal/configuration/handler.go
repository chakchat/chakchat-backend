package configuration

import "github.com/chakchat/chakchat-backend/messaging-service/internal/rest/handlers/chat"

type Handlers struct {
	PersonalChat       *chat.PersonalChatHandler
	GroupChat          *chat.GroupChatHandler
	GroupPhoto         *chat.GroupPhotoHandler
	SecretPersonalChat *chat.SecretPersonalChatHandler
	SecretGroup        *chat.SecretGroupHandler
	SecretGroupPhoto   *chat.SecretGroupPhotoHandler
	GenericChat        *chat.GenericChatHandler
}

func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		PersonalChat:       chat.NewPersonalChatHandler(services.PersonalChat),
		GroupChat:          chat.NewGroupChatHandler(services.GroupChat),
		GroupPhoto:         chat.NewGroupPhotoHandler(services.GroupPhoto),
		SecretPersonalChat: chat.NewSecretPersonalChatHandler(services.SecretPersonalChat),
		SecretGroup:        chat.NewSecretGroupHandler(services.SecretGroup),
		SecretGroupPhoto:   chat.NewSecretGroupPhotoHandler(services.SecretGroupPhoto),
		GenericChat:        chat.NewGenericChatHandler(services.GenericChat),
	}
}
