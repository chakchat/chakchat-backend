package dto

import (
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secpersonal"
	"github.com/google/uuid"
)

type SecretPersonalChatDTO struct {
	ID         uuid.UUID
	CreatedAt  int64
	Expiration *time.Duration
	Members    [2]uuid.UUID
}

func NewSecretPersonalChatDTO(c *secpersonal.SecretPersonalChat) SecretPersonalChatDTO {
	return SecretPersonalChatDTO{
		ID:         uuid.UUID(c.ID),
		CreatedAt:  int64(c.CreatedAt),
		Expiration: c.Exp,
		Members: [2]uuid.UUID{
			uuid.UUID(c.Members[0]),
			uuid.UUID(c.Members[1]),
		},
	}
}
