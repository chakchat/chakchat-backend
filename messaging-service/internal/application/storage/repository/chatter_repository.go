package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

// TODO: move it to GenericChatRepository
type ChatterRepository interface {
	FindChatter(context.Context, storage.ExecQuerier, domain.ChatID) (domain.Chatter, error)
}
