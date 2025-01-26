package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type PersonalChatDTO struct {
	ID      uuid.UUID
	Members [2]uuid.UUID

	Secret    bool
	Blocked   bool
	BlockedBy []uuid.UUID
	CreatedAt int64
}

func NewPersonalChatDTO(chat *domain.PersonalChat) PersonalChatDTO {
	blockedBy := make([]uuid.UUID, len(chat.BlockedBy))
	for i, u := range chat.BlockedBy {
		blockedBy[i] = uuid.UUID(u)
	}

	return PersonalChatDTO{
		ID: uuid.UUID(chat.ID),
		Members: [2]uuid.UUID{
			uuid.UUID(chat.Members[0]),
			uuid.UUID(chat.Members[1]),
		},
		Secret:    chat.Secret,
		Blocked:   chat.Blocked,
		BlockedBy: blockedBy,
		CreatedAt: int64(chat.CreatedAt),
	}
}
