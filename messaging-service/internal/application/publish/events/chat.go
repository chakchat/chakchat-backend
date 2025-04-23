package events

import (
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/google/uuid"
)

type ChatCreated struct {
	SenderID uuid.UUID    `json:"sender_id"`
	Chat     generic.Chat `json:"chat"`
}

type ChatDeleted struct {
	SenderID uuid.UUID `json:"sender_id"`
	ChatID   uuid.UUID `json:"chat_id"`
}

type ChatBlocked struct {
	SenderID uuid.UUID `json:"sender_id"`
	ChatID   uuid.UUID `json:"chat_id"`
}

type ChatUnblocked struct {
	SenderID uuid.UUID `json:"sender_id"`
	ChatID   uuid.UUID `json:"chat_id"`
}

type ExpirationSet struct {
	ChatID     uuid.UUID      `json:"chat_id"`
	SenderID   uuid.UUID      `json:"sender_id"`
	Expiration *time.Duration `json:"expiration"`
}

type GroupInfoUpdated struct {
	SenderID    uuid.UUID `json:"sender_id"`
	ChatID      uuid.UUID `json:"chat_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupPhoto  string    `json:"group_photo"`
}

type GroupMemberAdded struct {
	SenderID uuid.UUID   `json:"sender_id"`
	ChatID   uuid.UUID   `json:"chat_id"`
	Members  []uuid.UUID `json:"members"`
}

type GroupMembersRemoved struct {
	SenderID uuid.UUID   `json:"sender_id"`
	ChatID   uuid.UUID   `json:"chat_id"`
	Members  []uuid.UUID `json:"members"`
}
