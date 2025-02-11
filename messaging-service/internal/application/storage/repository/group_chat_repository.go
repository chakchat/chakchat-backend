package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/group"
)

type GroupChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, domain.ChatID) (*group.GroupChat, error)
	Update(context.Context, *group.GroupChat) (*group.GroupChat, error)
	Create(context.Context, *group.GroupChat) (*group.GroupChat, error)
	Delete(context.Context, domain.ChatID) error
}
