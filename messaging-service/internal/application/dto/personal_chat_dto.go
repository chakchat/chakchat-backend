package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/personal"
	"github.com/google/uuid"
)

type PersonalChatDTO struct {
	ID      uuid.UUID
	Members [2]uuid.UUID

	Blocked   bool
	BlockedBy []uuid.UUID
	CreatedAt int64
}

func NewPersonalChatDTO(chat *personal.PersonalChat) PersonalChatDTO {
	return PersonalChatDTO{
		ID: uuid.UUID(chat.ID),
		Members: [2]uuid.UUID{
			uuid.UUID(chat.Members[0]),
			uuid.UUID(chat.Members[1]),
		},
		BlockedBy: UUIDs(chat.BlockedBy),
		CreatedAt: int64(chat.CreatedAt),
	}
}
