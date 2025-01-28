package domain

import (
	"errors"
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

type TextMessage struct {
	ID       UpdateID
	ChatID   ChatID
	SenderID UserID

	Text   string
	SentAt Timestamp
	Read   bool
}

func NewTextMessage(sender UserID, text string, sentAt Timestamp) (TextMessage, error) {
	if text == "" {
		return TextMessage{}, ErrEmptyMessage
	}
	if len(text) > MaxTextMessageLen {
		return TextMessage{}, ErrTextMessageTooLong
	}

	return TextMessage{
		SenderID: sender,
		Text:     text,
		SentAt:   sentAt,
	}, nil
}

func (m *TextMessage) AssignToPersonalChat(chat *PersonalChat) error {
	if !chat.IsMember(m.SenderID) {
		return ErrUserNotMember
	}
	if chat.Blocked {
		return ErrChatIsBlocked
	}

	m.ChatID = chat.ID
	return nil
}

func (m *TextMessage) AssignToGroupChat(group *GroupChat) error {
	if !group.IsMember(m.SenderID) {
		return ErrUserNotMember
	}

	m.ChatID = group.ID
	return nil
}
