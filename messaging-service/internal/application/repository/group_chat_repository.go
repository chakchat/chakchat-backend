package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GroupChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, domain.ChatID) (*domain.GroupChat, error)
	Update(context.Context, *domain.GroupChat) (*domain.GroupChat, error)
	Create(context.Context, *domain.GroupChat) (*domain.GroupChat, error)
	Delete(context.Context, domain.ChatID) error
}
