package configuration

import "github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/chat"

type Services struct {
	PersonalChat       *chat.PersonalChatService
	GroupChat          *chat.GroupChatService
	GroupPhoto         *chat.GroupPhotoService
	SecretPersonalChat *chat.SecretPersonalChatService
	SecretGroup        *chat.SecretGroupChatService
	SecretGroupPhoto   *chat.SecretGroupPhotoService
}

func NewServices(db *DB, external *External) *Services {
	return &Services{
		PersonalChat:       chat.NewPersonalChatService(db.PersonalChat, external.Publisher),
		GroupChat:          chat.NewGroupChatService(db.GroupChat, external.Publisher),
		GroupPhoto:         chat.NewGroupPhotoService(db.GroupChat, external.FileStorage, external.Publisher),
		SecretPersonalChat: chat.NewSecretPersonalChatService(db.SecretPersonalChat, external.Publisher),
		SecretGroup:        chat.NewSecretGroupChatService(db.SecretGroupChat, external.Publisher),
		SecretGroupPhoto:   chat.NewSecretGroupPhotoService(db.SecretGroupChat, external.FileStorage, external.Publisher),
	}
}
