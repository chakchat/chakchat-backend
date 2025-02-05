package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secpersonal"
)

type SecretPersonalChatRepository interface {
	// Should return ErrNotFound if entity is not found
	FindById(context.Context, domain.ChatID) (*secpersonal.SecretPersonalChat, error)
	// Should return ErrNotFound if entity is not found.
	// Members order should NOT affect the result
	FindByMembers(context.Context, [2]domain.UserID) (*secpersonal.SecretPersonalChat, error)
	Update(context.Context, *secpersonal.SecretPersonalChat) (*secpersonal.SecretPersonalChat, error)
	Create(context.Context, *secpersonal.SecretPersonalChat) (*secpersonal.SecretPersonalChat, error)
	Delete(context.Context, domain.ChatID) error
}
