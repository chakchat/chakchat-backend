package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type TextMessageDTO struct {
	ChatID   uuid.UUID
	UpdateID int64
	SenderID uuid.UUID

	Text   string
	Edited *TextMessageEditedDTO

	CreatedAt int64
}

func NewTextMessageDTO(m *domain.TextMessage) TextMessageDTO {
	var edited *TextMessageEditedDTO
	if m.Edited != nil {
		editedDto := NewTextMessageEditedDTO(m.Edited)
		edited = &editedDto
	}

	return TextMessageDTO{
		ChatID:    uuid.UUID(m.ChatID),
		UpdateID:  int64(m.UpdateID),
		SenderID:  uuid.UUID(m.SenderID),
		Text:      m.Text,
		Edited:    edited,
		CreatedAt: int64(m.CreatedAt),
	}
}
