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
	// Should return ErrNotFound if entity is not found.
	// Members order should NOT affect the result
	FindByMembers(members [2]domain.UserID) (*domain.PersonalChat, error)
	Update(*domain.PersonalChat) (*domain.PersonalChat, error)
	Create(chat *domain.PersonalChat) (*domain.PersonalChat, error)
	Delete(d domain.ChatID) error
}
