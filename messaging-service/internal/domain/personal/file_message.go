package personal

import (
	"errors"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

const (
	MaxFileSize = 1 << 30
)

var (
	ErrFileTooBig = errors.New("file is too big")
)

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

func (c *PersonalChat) NewFileMessage(sender domain.UserID, file *FileMeta, replyTo *Message) (FileMessage, error) {
	if err := c.validateCanSend(sender); err != nil {
		return FileMessage{}, err
	}

	if err := validateFile(file); err != nil {
		return FileMessage{}, err
	}

	if replyTo != nil {
		if err := c.validateCanReply(sender, replyTo); err != nil {
			return FileMessage{}, err
		}
	}

	return FileMessage{
		Message: Message{
			Update: domain.Update{
				ChatID:   c.ChatID,
				SenderID: sender,
			},
			ReplyTo: replyTo,
		},
		File: *file,
	}, nil
}

func (c *PersonalChat) DeleteFileMessage(sender domain.UserID, m *FileMessage, mode domain.DeleteMode) error {
	return c.deleteMessage(sender, &m.Update, mode)
}

func validateFile(file *FileMeta) error {
	if file.FileSize > MaxFileSize {
		return ErrFileTooBig
	}
	return nil
}
