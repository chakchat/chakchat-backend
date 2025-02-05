package personal

import (
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

const (
	MaxFileSize = 1 << 30
)

var (
	ErrFileTooBig = errors.New("file is too big")
)

func (c *PersonalChat) NewFileMessage(sender domain.UserID, file *domain.FileMeta, replyTo *domain.Message) (domain.FileMessage, error) {
	if err := c.validateCanSend(sender); err != nil {
		return domain.FileMessage{}, err
	}

	if err := validateFile(file); err != nil {
		return domain.FileMessage{}, err
	}

	if replyTo != nil {
		if err := c.validateCanReply(sender, replyTo); err != nil {
			return domain.FileMessage{}, err
		}
	}

	return domain.FileMessage{
		Message: domain.Message{
			Update: domain.Update{
				ChatID:   c.ChatID,
				SenderID: sender,
			},
			ReplyTo: replyTo,
		},
		File: *file,
	}, nil
}

func (c *PersonalChat) DeleteFileMessage(sender domain.UserID, m *domain.FileMessage, mode domain.DeleteMode) error {
	return c.deleteMessage(sender, &m.Update, mode)
}

func validateFile(file *domain.FileMeta) error {
	if file.FileSize > MaxFileSize {
		return ErrFileTooBig
	}
	return nil
}
