package configuration

import "github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/chat"

type Services struct {
	PersonalChat       *chat.PersonalChatService
	GroupChat          *chat.GroupChatService
	GroupPhoto         *chat.GroupPhotoService
	SecretPersonalChat *chat.SecretPersonalChatService
	SecretGroup        *chat.SecretGroupChatService
	SecretGroupPhoto   *chat.SecretGroupPhotoService
	GenericChat        *chat.GenericChatService
}

func NewServices(db *DB, external *External) *Services {
	return &Services{
		PersonalChat:       chat.NewPersonalChatService(db.TxProvider, db.PersonalChat, external.Publisher),
		GroupChat:          chat.NewGroupChatService(db.TxProvider, db.GroupChat, external.Publisher),
		GroupPhoto:         chat.NewGroupPhotoService(db.TxProvider, db.GroupChat, external.FileStorage, external.Publisher),
		SecretPersonalChat: chat.NewSecretPersonalChatService(db.TxProvider, db.SecretPersonalChat, external.Publisher),
		SecretGroup:        chat.NewSecretGroupChatService(db.TxProvider, db.SecretGroupChat, external.Publisher),
		SecretGroupPhoto:   chat.NewSecretGroupPhotoService(db.TxProvider, db.SecretGroupChat, external.FileStorage, external.Publisher),
		GenericChat:        chat.NewGenericChatService(db.TxProvider, &db.GenericChat),
	}
}
