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

	ErrMessageAlreadyDeleted = errors.New("message is already deleted")
)

const (
	MaxTextMessageLen = 1000
)

type DeletedMode int

const (
	DeletedModeNone = iota
	DeletedModeForSender
	DeletedModeForAll
)

type Message struct {
	ID       UpdateID
	ChatID   ChatID
	SenderID UserID

	SentAt  Timestamp
	Deleted DeletedMode
	Read    bool
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

type MessageDeleted struct {
	UpdateID  UpdateID
	ChatID    ChatID
	MessageID UpdateID

	DeletedMode DeletedMode
}

func (m *Message) Delete(onlyMe bool) (MessageDeleted, error) {
	if onlyMe {
		if m.Deleted != DeletedModeNone {
			return MessageDeleted{}, ErrMessageAlreadyDeleted
		}
		m.Deleted = DeletedModeForSender
	} else {
		if m.Deleted == DeletedModeForAll {
			return MessageDeleted{}, ErrMessageAlreadyDeleted
		}
		m.Deleted = DeletedModeForAll
	}

	return MessageDeleted{
		ChatID:      m.ChatID,
		MessageID:   m.ID,
		DeletedMode: m.Deleted,
	}, nil
}

type TextMessage struct {
	Message

	Text     string
	EditedAt Timestamp
}

func NewTextMessage(sender UserID, text string, sentAt Timestamp) (TextMessage, error) {
	if err := validateText(text); err != nil {
		return TextMessage{}, err
	}

	return TextMessage{
		Message: Message{
			SenderID: sender,
			SentAt:   sentAt,
		},
		Text: text,
	}, nil
}

func (m *TextMessage) Edit(newText string) error {
	if err := validateText(newText); err != nil {
		return err
	}

	m.Text = newText
	m.EditedAt = Timestamp(TimeFunc().Unix())
	return nil
}

func validateText(text string) error {
	if text == "" {
		return ErrEmptyMessage
	}
	if len(text) > MaxTextMessageLen {
		return ErrTextMessageTooLong
	}
	return nil
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
