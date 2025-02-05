package personal

import (
	"errors"
	"unicode/utf8"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

const (
	MaxTextRunesCount = 2000
)

var (
	ErrTooMuchTextRunes = errors.New("too much runes in text")
	ErrChatBlocked      = errors.New("chat is blocked")
	ErrTextEmpty        = errors.New("the text is empty")
)

type TextMessage struct {
	Message

	Text   string
	Edited *TextMessageEdited
}

type TextMessageEdited struct {
	domain.Update

	MessageID domain.UpdateID
}

func (c *PersonalChat) NewTextMessage(sender domain.UserID, text string, replyTo *Message) (TextMessage, error) {
	if err := c.validateCanSend(sender); err != nil {
		return TextMessage{}, err
	}

	if replyTo != nil {
		if err := c.validateCanReply(sender, replyTo); err != nil {
			return TextMessage{}, err
		}
	}

	if err := validateText(text); err != nil {
		return TextMessage{}, err
	}

	return TextMessage{
		Message: Message{
			Update: domain.Update{
				ChatID:   c.ChatID,
				SenderID: sender,
			},
			ReplyTo: replyTo,
		},
		Text: text,
	}, nil
}

func (c *PersonalChat) EditTextMessage(sender domain.UserID, m *TextMessage, newText string) error {
	if err := c.validateCanSend(sender); err != nil {
		return err
	}

	if c.ChatID != m.ChatID {
		return domain.ErrUpdateNotFromChat
	}

	if m.SenderID != sender {
		return domain.ErrUserNotSender
	}

	if m.DeletedFor(sender) {
		return domain.ErrUpdateDeleted
	}

	if err := validateText(newText); err != nil {
		return err
	}

	m.Text = newText
	m.Edited = &TextMessageEdited{
		Update: domain.Update{
			ChatID:   c.ChatID,
			SenderID: sender,
		},
		MessageID: m.UpdateID,
	}
	return nil
}

func (c *PersonalChat) DeleteTextMessage(sender domain.UserID, m *TextMessage, mode domain.DeleteMode) error {
	return c.deleteMessage(sender, &m.Update, mode)
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
