package request

import "github.com/google/uuid"

type SendTextMessage struct {
	ChatID         uuid.UUID
	SenderID       uuid.UUID
	Text           string
	ReplyToMessage *int64
}

type EditTextMessage struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	MessageID int64
	NewText   string
}

type DeleteMessage struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	MessageID  int64
	DeleteMode string
}
