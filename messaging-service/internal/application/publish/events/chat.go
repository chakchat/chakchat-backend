package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	ChatTypePersonal       = "personal"
	ChatTypeSecretPersonal = "secret_personal"
	ChatTypeGroup          = "group"
	ChatTypeSecretGroup    = "secret_group"
)

type ChatCreated struct {
	ChatID   uuid.UUID `json:"chat_id"`
	ChatType string    `json:"chat_type"`
}

type ChatDeleted struct {
	ChatID uuid.UUID `json:"chat_id"`
}

type ExpirationSet struct {
	ChatID     uuid.UUID      `json:"chat_id"`
	SenderID   uuid.UUID      `json:"sender_id"`
	Expiration *time.Duration `json:"duration"`
}

type GroupInfoUpdated struct {
	ChatID        uuid.UUID `json:"chat_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	GroupPhotoURL string    `json:"group_photo_url"`
}

type GroupMemberAdded struct {
	ChatID   uuid.UUID `json:"chat_id"`
	MemberID uuid.UUID `json:"member_id"`
}

type GroupMemberRemoved struct {
	ChatID   uuid.UUID `json:"chat_id"`
	MemberID uuid.UUID `json:"member_id"`
}

type ChatBlocked struct {
	ChatID uuid.UUID `json:"chat_id"`
}

type ChatUnblocked struct {
	ChatID uuid.UUID `json:"chat_id"`
}
