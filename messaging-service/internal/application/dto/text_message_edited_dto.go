package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type TextMessageEditedDTO struct {
	ChatID   uuid.UUID
	UpdateID int64
	SenderID uuid.UUID

	MessageID int64
	NewText   string

	CreatedAt int64
}

func NewTextMessageEditedDTO(dom *domain.TextMessageEdited) TextMessageEditedDTO {
	return TextMessageEditedDTO{
		ChatID:    uuid.UUID(dom.ChatID),
		UpdateID:  int64(dom.UpdateID),
		SenderID:  uuid.UUID(dom.SenderID),
		MessageID: int64(dom.MessageID),
		NewText:   dom.NewText,
		CreatedAt: int64(dom.CreatedAt),
	}
}
