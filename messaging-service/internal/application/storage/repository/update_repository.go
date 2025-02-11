package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type UpdateRepository interface {
	FindGenericMessage(context.Context, domain.UpdateID) (*domain.Message, error)
	CreateTextMessage(context.Context, *domain.TextMessage) (*domain.TextMessage, error)
	FindTextMessage(context.Context, domain.UpdateID) (*domain.TextMessage, error)
	UpdateTextMessage(context.Context, *domain.TextMessage) (*domain.TextMessage, error)
	CreateTextMessageEdited(context.Context, *domain.TextMessageEdited) (*domain.TextMessageEdited, error)
}
