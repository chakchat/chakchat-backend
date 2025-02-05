package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Update

	ReplyTo *Message
}

type TextMessage struct {
	Message

	Text   string
	Edited *TextMessageEdited
}

type TextMessageEdited struct {
	Update

	MessageID UpdateID
}

type FileMeta struct {
	FileId    uuid.UUID
	FileName  string
	MimeType  string
	FileSize  int64
	FileUrl   string
	CreatedAt time.Time
}

type FileMessage struct {
	Message
	File FileMeta
}
