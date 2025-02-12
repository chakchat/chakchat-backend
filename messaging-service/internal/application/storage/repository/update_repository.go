package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type UpdateRepository interface {
	FindGenericMessage(context.Context, domain.ChatID, domain.UpdateID) (*domain.Message, error)
	CreateTextMessage(context.Context, *domain.TextMessage) (*domain.TextMessage, error)
	FindTextMessage(context.Context, domain.ChatID, domain.UpdateID) (*domain.TextMessage, error)
	UpdateTextMessage(context.Context, *domain.TextMessage) (*domain.TextMessage, error)
	CreateTextMessageEdited(context.Context, *domain.TextMessageEdited) (*domain.TextMessageEdited, error)
	DeleteMessage(context.Context, domain.ChatID, domain.UpdateID) error
	CreateUpdateDeleted(context.Context, *domain.UpdateDeleted) (*domain.UpdateDeleted, error)
	CreateReaction(context.Context, *domain.Reaction) (*domain.Reaction, error)
}
