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

type UpdateDeleted struct {
	ChatID   uuid.UUID `json:"chat_id"`
	UpdateID int64     `json:"update_id"`
	SenderID uuid.UUID `json:"sender_id"`

	DeletedID  int64  `json:"deleted_id"`
	DeleteMode string `json:"delete_mode"`

	CreatedAt int64 `json:"created_at"`
}

type ReactionSent struct {
	ChatID   uuid.UUID `json:"chat_id"`
	UpdateID int64     `json:"update_id"`
	SenderID uuid.UUID `json:"sender_id"`

	CreatedAt    int64  `json:"created_at"`
	ReactionType string `json:"reaction_type"`
}
