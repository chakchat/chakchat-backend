package domain

import (
	"errors"

	"github.com/google/uuid"
)

type (
	ChatID uuid.UUID
	UserID uuid.UUID
)

func NewUserID(id string) (UserID, error) {
	userId, err := uuid.Parse(id)
	return UserID(userId), err
}

func NewChatID() ChatID {
	return ChatID(uuid.New())
}

var (
	ErrUserNotMember = errors.New("user is not member of a chat")
)

type Chat struct {
	ChatID    ChatID
	CreatedAt Timestamp
}
