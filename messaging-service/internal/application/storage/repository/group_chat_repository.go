package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/group"
)

type GroupChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, storage.ExecQuerier, domain.ChatID) (*group.GroupChat, error)
	Update(context.Context, storage.ExecQuerier, *group.GroupChat) (*group.GroupChat, error)
	Create(context.Context, storage.ExecQuerier, *group.GroupChat) (*group.GroupChat, error)
	Delete(context.Context, storage.ExecQuerier, domain.ChatID) error
}
