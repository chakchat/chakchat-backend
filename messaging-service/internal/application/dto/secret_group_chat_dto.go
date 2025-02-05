package dto

import (
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secgroup"
	"github.com/google/uuid"
)

type SecretGroupChatDTO struct {
	ID        uuid.UUID
	CreatedAt int64

	Admin   uuid.UUID
	Members []uuid.UUID

	Name          string
	Description   string
	GroupPhotoURL string

	Expiration *time.Duration
}

func NewSecretGroupChatDTO(g *secgroup.SecretGroupChat) SecretGroupChatDTO {
	return SecretGroupChatDTO{
		ID:            uuid.UUID(g.ID),
		CreatedAt:     int64(g.CreatedAt),
		Admin:         uuid.UUID(g.Admin),
		Members:       UUIDs(g.Members),
		Name:          g.Name,
		Description:   g.Description,
		GroupPhotoURL: string(g.GroupPhoto),
		Expiration:    g.Expiration,
	}
}
