package events

import "github.com/google/uuid"

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

type GroupInfoUpdated struct {
	ChatID        uuid.UUID `json:"chat_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	GroupPhotoURL string    `json:"group_photo_url"`
}
