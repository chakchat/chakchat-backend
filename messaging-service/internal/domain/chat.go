package domain

import (
	"errors"

	"github.com/google/uuid"
)

type (
	ChatID uuid.UUID
	UserID uuid.UUID
)

var (
	ErrUserNotMember = errors.New("user is not member of a chat")
)

type Chat struct {
	ChatID    ChatID
	CreatedAt Timestamp
}

func NewChatID() ChatID {
	return ChatID(uuid.New())
}
