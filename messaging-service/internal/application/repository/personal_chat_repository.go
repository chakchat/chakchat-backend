package repository

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

var (
	ErrNotFound = errors.New("not found")
)

//go:generate mockery
type PersonalChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, domain.ChatID) (*domain.PersonalChat, error)
	// Should return ErrNotFound if entity is not found.
	// Members order should NOT affect the result
	FindByMembers(context.Context, [2]domain.UserID) (*domain.PersonalChat, error)
	Update(context.Context, *domain.PersonalChat) (*domain.PersonalChat, error)
	Create(context.Context, *domain.PersonalChat) (*domain.PersonalChat, error)
	Delete(context.Context, domain.ChatID) error
}
