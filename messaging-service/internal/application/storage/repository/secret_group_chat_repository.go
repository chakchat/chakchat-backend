package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secgroup"
)

type SecretGroupChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, storage.ExecQuerier, domain.ChatID) (*secgroup.SecretGroupChat, error)
	Update(context.Context, storage.ExecQuerier, *secgroup.SecretGroupChat) (*secgroup.SecretGroupChat, error)
	Create(context.Context, storage.ExecQuerier, *secgroup.SecretGroupChat) (*secgroup.SecretGroupChat, error)
	Delete(context.Context, storage.ExecQuerier, domain.ChatID) error
}
