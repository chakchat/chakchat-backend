package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type ChatterRepository interface {
	FindChatter(context.Context, domain.ChatID) (domain.Chatter, error)
}
