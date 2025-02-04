package domain

import "github.com/google/uuid"

type (
	ChatID uuid.UUID
	UserID uuid.UUID
)

type Chat struct {
	ChatID    ChatID
	CreatedAt Timestamp
}

func NewChatID() ChatID {
	return ChatID(uuid.New())
}
