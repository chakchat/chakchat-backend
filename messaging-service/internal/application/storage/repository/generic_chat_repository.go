package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GenericChatRepository interface {
	// Should return empty slice if not found
	GetByMemberID(context.Context, domain.UserID) ([]services.GenericChat, error)
	// Should return ErrNotFound if not found
	GetByChatID(context.Context, domain.ChatID) (*services.GenericChat, error)
	// Should return ErrNotFound if not found
	GetChatType(context.Context, domain.ChatID) (string, error)
}
