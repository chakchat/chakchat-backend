package domain

import (
	"unicode/utf8"
)

const (
	MaxTextRunesCount = 2000
)

type TextMessage struct {
	Message

	Text   string
	Edited *TextMessageEdited
}

type TextMessageEdited struct {
	Update

	MessageID UpdateID
	NewText   string
}

func NewTextMessage(chat Chatter, sender UserID, text string, replyTo *Message) (*TextMessage, error) {
	if err := chat.ValidateCanSend(sender); err != nil {
		return nil, err
	}

	if replyTo != nil {
		if err := validateCanReply(chat, sender, replyTo); err != nil {
			return nil, err
		}
	}

	if err := validateText(text); err != nil {
		return nil, err
	}

	var replyToID *UpdateID
	if replyTo != nil {
		replyToID = &replyTo.UpdateID
	}

	return &TextMessage{
		Message: Message{
			Update: Update{
				ChatID:   chat.ChatID(),
				SenderID: sender,
			},
			ReplyTo: replyToID,
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
		Update:    Update{ChatID: chat.ChatID(), SenderID: sender},
		MessageID: m.UpdateID,
		NewText:   newText,
	}
	return nil
}

func (m *TextMessage) Forward(fromChat Chatter, sender UserID, toChat Chatter) (*TextMessage, error) {
	if !fromChat.IsMember(sender) {
		return nil, ErrUserNotMember
	}
	if err := toChat.ValidateCanSend(sender); err != nil {
		return nil, err
	}

	return &TextMessage{
		Message: Message{
			Update: Update{
				ChatID:   toChat.ChatID(),
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
		return ErrTooManyTextRunes
	}
	return nil
}
