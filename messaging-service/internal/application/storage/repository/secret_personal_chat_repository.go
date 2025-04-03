package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secpersonal"
)

type SecretPersonalChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, storage.ExecQuerier, domain.ChatID) (*secpersonal.SecretPersonalChat, error)
	// Should return ErrNotFound if entity is not found.
	// Members order should NOT affect the result
	FindByMembers(context.Context, storage.ExecQuerier, [2]domain.UserID) (*secpersonal.SecretPersonalChat, error)
	Update(context.Context, storage.ExecQuerier, *secpersonal.SecretPersonalChat) (*secpersonal.SecretPersonalChat, error)
	Create(context.Context, storage.ExecQuerier, *secpersonal.SecretPersonalChat) (*secpersonal.SecretPersonalChat, error)
	Delete(context.Context, storage.ExecQuerier, domain.ChatID) error
}
