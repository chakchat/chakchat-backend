package configuration

import "github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/chat"

type Services struct {
	PersonalChat *chat.PersonalChatService
}

func NewService(db *DB, external *External) *Services {
	return &Services{
		PersonalChat: chat.NewPersonalChatService(db.PersonalChat, external.Publisher),
	}
}
