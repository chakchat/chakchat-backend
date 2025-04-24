package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type ReactionDTO struct {
	UpdateID int64
	ChatID   uuid.UUID
	SenderID uuid.UUID

	CreatedAt    int64
	MessageID    int64
	ReactionType string
}

func NewReactionDTO(r *domain.Reaction) ReactionDTO {
	return ReactionDTO{
		UpdateID:     int64(r.UpdateID),
		ChatID:       uuid.UUID(r.ChatID),
		SenderID:     uuid.UUID(r.SenderID),
		CreatedAt:    int64(r.CreatedAt),
		MessageID:    int64(r.MessageID),
		ReactionType: string(r.Type),
	}
}
