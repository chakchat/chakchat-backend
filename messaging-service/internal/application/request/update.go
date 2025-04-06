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

type SendReaction struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	MessageID    int64
	ReactionType string
}

type DeleteReaction struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	ReactionID int64
}

type ForwardMessage struct {
	ToChatID uuid.UUID
	SenderID uuid.UUID

	MessageID  int64
	FromChatID uuid.UUID
}

type SendFileMessage struct {
	ChatID         uuid.UUID
	SenderID       uuid.UUID
	FileID         uuid.UUID
	ReplyToMessage *int64
}

type SendSecretUpdate struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	Payload              []byte
	InitializationVector []byte
	KeyHash              []byte
}

type DeleteSecretUpdate struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	SecretUpdateID int64
}

type GetUpdatesRange struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	From, To int64
}

type GetUpdate struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID

	UpdateID int64
}
