package domain

import (
	"errors"

	"github.com/google/uuid"
)

type UpdateID int64

var (
	ErrEmptyMessage       = errors.New("empty message")
	ErrTextMessageTooLong = errors.New("text message is too long")
	ErrChatIsBlocked      = errors.New("chat is blocked")
)

const (
	MaxTextMessageLen = 1000
)

type Message struct {
	ID       UpdateID
	ChatID   ChatID
	SenderID UserID

	SentAt Timestamp
	Read   bool
}

func (m *Message) AssignToPersonalChat(chat *PersonalChat) error {
	if !chat.IsMember(m.SenderID) {
		return ErrUserNotMember
	}
	if chat.Blocked {
		return ErrChatIsBlocked
	}

	m.ChatID = chat.ID
	return nil
}

func (m *Message) AssignToGroupChat(group *GroupChat) error {
	if !group.IsMember(m.SenderID) {
		return ErrUserNotMember
	}

	m.ChatID = group.ID
	return nil
}

type TextMessage struct {
	Message
	Text string
}

func NewTextMessage(sender UserID, text string, sentAt Timestamp) (TextMessage, error) {
	if text == "" {
		return TextMessage{}, ErrEmptyMessage
	}
	if len(text) > MaxTextMessageLen {
		return TextMessage{}, ErrTextMessageTooLong
	}

	return TextMessage{
		Message: Message{
			SenderID: sender,
			SentAt:   sentAt,
		},
		Text: text,
	}, nil
}

type FileMeta struct {
	FileId    uuid.UUID
	FileName  string
	MimeType  string
	FileSize  int64
	FileUrl   URL
	CreatedAt Timestamp
}

type FileMessage struct {
	Message
	File FileMeta
}

func NewFileMessage(sender UserID, file FileMeta, sentAt Timestamp) (FileMessage, error) {
	return FileMessage{
		Message: Message{
			SenderID: sender,
			SentAt:   sentAt,
		},
		File: file,
	}, nil
}
