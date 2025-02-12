package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type UpdateRepository interface {
	FindGenericMessage(context.Context, domain.ChatID, domain.UpdateID) (*domain.Message, error)
	DeleteMessage(context.Context, domain.ChatID, domain.UpdateID) error
	CreateUpdateDeleted(context.Context, *domain.UpdateDeleted) (*domain.UpdateDeleted, error)

	CreateTextMessage(context.Context, *domain.TextMessage) (*domain.TextMessage, error)
	CreateTextMessageEdited(context.Context, *domain.TextMessageEdited) (*domain.TextMessageEdited, error)
	FindTextMessage(context.Context, domain.ChatID, domain.UpdateID) (*domain.TextMessage, error)
	UpdateTextMessage(context.Context, *domain.TextMessage) (*domain.TextMessage, error)

	CreateReaction(context.Context, *domain.Reaction) (*domain.Reaction, error)
	FindReaction(context.Context, domain.ChatID, domain.UpdateID) (*domain.Reaction, error)
	DeleteReaction(context.Context, domain.ChatID, domain.UpdateID) error
}
