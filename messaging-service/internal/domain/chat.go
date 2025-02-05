package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrChatBlocked = errors.New("chat is blocked")
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
	ID        ChatID
	CreatedAt Timestamp
}

func (c *Chat) ChatID() ChatID {
	return c.ID
}

type Chatter interface {
	ChatID() ChatID
	ValidateCanSend(UserID) error
}
