package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GenericChatRepository interface {
	// Should return empty slice if not found
	GetByMemberID(context.Context, storage.ExecQuerier, domain.UserID) ([]generic.Chat, error)
	// Should return ErrNotFound if not found
	GetByChatID(context.Context, storage.ExecQuerier, domain.ChatID) (*generic.Chat, error)
	// Should return ErrNotFound if not found
	GetChatType(context.Context, storage.ExecQuerier, domain.ChatID) (string, error)
}
