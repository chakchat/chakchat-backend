package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type UpdateDeletedDTO struct {
	ChatID   uuid.UUID
	UpdateID int64
	SenderID uuid.UUID

	DeletedID  int64
	DeleteMode string

	CreatedAt int64
}

func NewUpdateDeletedDTO(dom *domain.UpdateDeleted) UpdateDeletedDTO {
	return UpdateDeletedDTO{
		ChatID:     uuid.UUID(dom.ChatID),
		UpdateID:   int64(dom.UpdateID),
		SenderID:   uuid.UUID(dom.SenderID),
		DeletedID:  int64(dom.DeletedID),
		DeleteMode: string(dom.Mode),
		CreatedAt:  int64(dom.CreatedAt),
	}
}
