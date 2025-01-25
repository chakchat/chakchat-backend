package repository

import (
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

type PersonalChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(chatId domain.ChatID) (*domain.PersonalChat, error)
	Update(*domain.PersonalChat) error
}
