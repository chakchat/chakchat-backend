package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GenericUpdateRepository interface {
	GetRange(
		ctx context.Context,
		db storage.ExecQuerier,
		visibleTo domain.UserID,
		chatID domain.ChatID,
		from, to domain.UpdateID,
	) ([]services.GenericUpdate, error)
	Get(
		ctx context.Context,
		db storage.ExecQuerier,
		visibleTo domain.UserID,
		chatID domain.ChatID,
		updateID domain.UpdateID,
	) (*services.GenericUpdate, error)
}
