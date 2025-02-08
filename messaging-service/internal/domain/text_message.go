package domain

import (
	"errors"
	"unicode/utf8"
)

const (
	MaxTextRunesCount = 2000
)

var (
	ErrTooMuchTextRunes = errors.New("too much runes in text")
	ErrTextEmpty        = errors.New("the text is empty")
)

type TextMessage struct {
	Message

	Text   string
	Edited *TextMessageEdited
}

type TextMessageEdited struct {
	Update

	MessageID UpdateID
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

func (m *TextMessage) Forward(chat Chatter, sender UserID, destChat Chatter) (TextMessage, error) {
	if !chat.IsMember(sender) {
		return TextMessage{}, ErrUserNotMember
	}
	if err := destChat.ValidateCanSend(sender); err != nil {
		return TextMessage{}, err
	}

	return TextMessage{
		Message: Message{
			Update: Update{
				ChatID:   destChat.ChatID(),
				SenderID: sender,
			},
			Forwarded: true,
		},
		Text: m.Text,
	}, nil
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
