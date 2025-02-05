package events

import "github.com/google/uuid"

type ChatCreatedEvent struct {
	ChatID   uuid.UUID `json:"chat_id"`
	ChatType string    `json:"chat_type"`
}
