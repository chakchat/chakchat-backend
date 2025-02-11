package events

import "github.com/google/uuid"

type TextMessageSent struct {
	ChatID   uuid.UUID `json:"chat_id"`
	UpdateID int64     `json:"update_id"`
	SenderID uuid.UUID `json:"sender_id"`

	Text      string `json:"text"`
	CreatedAt int64  `json:"created_at"`
}

type TextMessageEdited struct {
	ChatID   uuid.UUID `json:"chat_id"`
	UpdateID int64     `json:"update_id"`
	SenderID uuid.UUID `json:"sender_id"`

	MessageID int64  `json:"message_id"`
	NewText   string `json:"text"`
	CreatedAt int64  `json:"created_at"`
}
