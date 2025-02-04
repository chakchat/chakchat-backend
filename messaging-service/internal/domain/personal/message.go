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
)

type TextMessage struct {
	domain.Update

	Text string
}

func (c *PersonalChat) NewTextMessage(sender domain.UserID, text string) (TextMessage, error) {
	if err := c.validateCanSend(sender); err != nil {
		return TextMessage{}, err
	}

	if utf8.RuneCountInString(text) > MaxTextRunesCount {
		return TextMessage{}, ErrTooMuchTextRunes
	}

	return TextMessage{
		Update: domain.Update{
			ChatID:   c.ChatID,
			SenderID: sender,
		},
		Text: text,
	}, nil
}

func (c *PersonalChat) validateCanSend(sender domain.UserID) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	if c.Blocked() {
		return ErrChatBlocked
	}
	return nil
}
