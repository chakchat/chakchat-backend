package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type UpdateRepository interface {
	FindGenericMessage(context.Context, domain.UpdateID) (*domain.Message, error)
	CreateTextMessage(context.Context, *domain.TextMessage) (*domain.TextMessage, error)
}
