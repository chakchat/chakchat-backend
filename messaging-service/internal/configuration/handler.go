package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/handlers/chat"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/handlers/update"
)

type Handlers struct {
	PersonalChat       *chat.PersonalChatHandler
	GroupChat          *chat.GroupChatHandler
	GroupPhoto         *chat.GroupPhotoHandler
	SecretPersonalChat *chat.SecretPersonalChatHandler
	SecretGroup        *chat.SecretGroupHandler
	SecretGroupPhoto   *chat.SecretGroupPhotoHandler
	GenericChat        *chat.GenericChatHandler

	PersonalUpdate       *update.PersonalUpdateHandler
	PersonalFile         *update.PersonalFileHandler
	GroupUpdate          *update.GroupUpdateHandler
	GroupFile            *update.GroupFileHandler
	SecretPersonalUpdate *update.SecretPersonalUpdateHandler
	SecretGroupUpdate    *update.SecretGroupUpdateHandler
}

func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		PersonalChat:         chat.NewPersonalChatHandler(services.PersonalChat),
		GroupChat:            chat.NewGroupChatHandler(services.GroupChat),
		GroupPhoto:           chat.NewGroupPhotoHandler(services.GroupPhoto),
		SecretPersonalChat:   chat.NewSecretPersonalChatHandler(services.SecretPersonalChat),
		SecretGroup:          chat.NewSecretGroupHandler(services.SecretGroup),
		SecretGroupPhoto:     chat.NewSecretGroupPhotoHandler(services.SecretGroupPhoto),
		GenericChat:          chat.NewGenericChatHandler(services.GenericChat),
		PersonalUpdate:       update.NewPersonalUpdateHandler(services.GroupUpdate),
		PersonalFile:         update.NewFileHandler(services.GroupFile),
		GroupUpdate:          update.NewGroupUpdateHandler(services.GroupUpdate),
		GroupFile:            update.NewGroupFileHandler(services.GroupFile),
		SecretPersonalUpdate: update.NewSecretPersonalUpdateHandler(services.SecretGroupUpdate),
		SecretGroupUpdate:    update.NewSecretGroupUpdateHandler(services.SecretGroupUpdate),
	}
}
