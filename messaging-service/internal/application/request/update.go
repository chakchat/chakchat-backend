package request

import "github.com/google/uuid"

type SendTextMessage struct {
	ChatID         uuid.UUID
	SenderID       uuid.UUID
	Text           string
	ReplyToMessage *int64
}
