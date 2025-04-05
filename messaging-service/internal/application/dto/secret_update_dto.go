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
	KeyHash              string

	CreatedAt int64
}

func NewSecretUpdateDTO(dom *domain.SecretUpdate) SecretUpdateDTO {
	return SecretUpdateDTO{
		ChatID:               uuid.UUID(dom.ChatID),
		UpdateID:             int64(dom.UpdateID),
		SenderID:             uuid.UUID(dom.SenderID),
		Payload:              dom.Data.Payload,
		InitializationVector: dom.Data.IV,
		KeyHash:              string(dom.Data.KeyHash),
		CreatedAt:            int64(dom.CreatedAt),
	}
}
