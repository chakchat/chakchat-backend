package configuration

import "github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/chat"

type Service struct {
	PersonalChatService *chat.PersonalChatService
}

func NewService(db *DB, external *External) *Service {
	return &Service{
		PersonalChatService: chat.NewPersonalChatService(db.PersonalChatRepo, external.Publisher),
	}
}
