package domain

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MaxTextRunesCount = 2000
	MaxFileSize       = 1 << 30
)

var (
	ErrFileTooBig       = errors.New("file is too big")
	ErrTooMuchTextRunes = errors.New("too much runes in text")
	ErrChatBlocked      = errors.New("chat is blocked")
	ErrTextEmpty        = errors.New("the text is empty")
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

func (m *TextMessage) EditTextMessage(chat Chatter, sender UserID, newText string) error {
	if err := chat.ValidateCanSend(sender); err != nil {
		return err
	}

	if chat.ChatID() != m.ChatID {
		return ErrUpdateNotFromChat
	}

	if m.SenderID != sender {
		return ErrUserNotSender
	}

	if m.DeletedFor(sender) {
		return ErrUpdateDeleted
	}

	if err := validateText(newText); err != nil {
		return err
	}

	m.Text = newText
	m.Edited = &TextMessageEdited{
		Update: Update{
			ChatID:   chat.ChatID(),
			SenderID: sender,
		},
		MessageID: m.UpdateID,
	}
	return nil
}

func NewTextMessage(chat Chatter, sender UserID, text string, replyTo *Message) (TextMessage, error) {
	if err := chat.ValidateCanSend(sender); err != nil {
		return TextMessage{}, err
	}

	if replyTo != nil {
		if err := validateCanReply(chat, sender, replyTo); err != nil {
			return TextMessage{}, err
		}
	}

	if err := validateText(text); err != nil {
		return TextMessage{}, err
	}

	return TextMessage{
		Message: Message{
			Update: Update{
				ChatID:   chat.ChatID(),
				SenderID: sender,
			},
			ReplyTo: replyTo,
		},
		Text: text,
	}, nil
}

func (m *TextMessage) Edit(chat Chatter, sender UserID, newText string) error {
	if err := chat.ValidateCanSend(sender); err != nil {
		return err
	}

	if chat.ChatID() != m.ChatID {
		return ErrUpdateNotFromChat
	}

	if m.SenderID != sender {
		return ErrUserNotSender
	}

	if m.DeletedFor(sender) {
		return ErrUpdateDeleted
	}

	if err := validateText(newText); err != nil {
		return err
	}

	m.Text = newText
	m.Edited = &TextMessageEdited{
		Update: Update{
			ChatID:   chat.ChatID(),
			SenderID: sender,
		},
		MessageID: m.UpdateID,
	}
	return nil
}

func NewFileMessage(chat Chatter, sender UserID, file *FileMeta, replyTo *Message) (FileMessage, error) {
	if err := chat.ValidateCanSend(sender); err != nil {
		return FileMessage{}, err
	}

	if err := validateFile(file); err != nil {
		return FileMessage{}, err
	}

	if replyTo != nil {
		if err := validateCanReply(chat, sender, replyTo); err != nil {
			return FileMessage{}, err
		}
	}

	return FileMessage{
		Message: Message{
			Update: Update{
				ChatID:   chat.ChatID(),
				SenderID: sender,
			},
			ReplyTo: replyTo,
		},
		File: *file,
	}, nil
}

func (m *Message) Delete(chat Chatter, sender UserID, mode DeleteMode) error {
	if err := chat.ValidateCanSend(sender); err != nil {
		return err
	}

	if chat.ChatID() != m.ChatID {
		return ErrUpdateNotFromChat
	}

	if m.DeletedFor(sender) {
		return ErrUpdateDeleted
	}

	m.AddDeletion(sender, mode)
	return nil
}

func validateFile(file *FileMeta) error {
	if file.FileSize > MaxFileSize {
		return ErrFileTooBig
	}
	return nil
}

func validateText(text string) error {
	if text == "" {
		return ErrTextEmpty
	}
	if utf8.RuneCountInString(text) > MaxTextRunesCount {
		return ErrTooMuchTextRunes
	}
	return nil
}

func validateCanReply(chat Chatter, sender UserID, replyTo *Message) error {
	if replyTo.DeletedFor(sender) {
		return ErrUpdateDeleted
	}
	if chat.ChatID() != replyTo.ChatID {
		return ErrUpdateNotFromChat
	}
	return nil
}
