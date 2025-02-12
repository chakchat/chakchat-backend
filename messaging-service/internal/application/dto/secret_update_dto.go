package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type SecretUpdateDTO struct {
	ChatID   uuid.UUID
	UpdateID int64
	SenderID uuid.UUID

	Payload              []byte
	InitializationVector []byte
	KeyID                uuid.UUID

	CreatedAt int64
}

func NewSecretUpdateDTO(dom *domain.SecretUpdate) SecretUpdateDTO {
	return SecretUpdateDTO{
		ChatID:               uuid.UUID(dom.ChatID),
		UpdateID:             int64(dom.UpdateID),
		SenderID:             uuid.UUID(dom.SenderID),
		Payload:              dom.Data.Payload,
		InitializationVector: dom.Data.IV,
		KeyID:                uuid.UUID(dom.Data.KeyID),
		CreatedAt:            int64(dom.CreatedAt),
	}
}
