package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type UpdateRepository interface {
	FindGenericMessage(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.Message, error)
	DeleteUpdate(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) error
	CreateUpdateDeleted(context.Context, storage.ExecQuerier, *domain.UpdateDeleted) (*domain.UpdateDeleted, error)

	CreateTextMessage(context.Context, storage.ExecQuerier, *domain.TextMessage) (*domain.TextMessage, error)
	CreateTextMessageEdited(context.Context, storage.ExecQuerier, *domain.TextMessageEdited) (*domain.TextMessageEdited, error)
	FindTextMessage(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.TextMessage, error)
	UpdateTextMessage(context.Context, storage.ExecQuerier, *domain.TextMessage) (*domain.TextMessage, error)

	CreateReaction(context.Context, storage.ExecQuerier, *domain.Reaction) (*domain.Reaction, error)
	FindReaction(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.Reaction, error)

	FindFileMessage(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.FileMessage, error)
	CreateFileMessage(context.Context, storage.ExecQuerier, *domain.FileMessage) (*domain.FileMessage, error)
}

type SecretUpdateRepository interface {
	CreateSecretUpdate(context.Context, storage.ExecQuerier, *domain.SecretUpdate) (*domain.SecretUpdate, error)
	FindSecretUpdate(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) (*domain.SecretUpdate, error)
	DeleteSecretUpdate(context.Context, storage.ExecQuerier, domain.ChatID, domain.UpdateID) error
	CreateUpdateDeleted(context.Context, storage.ExecQuerier, *domain.UpdateDeleted) (*domain.UpdateDeleted, error)
}
