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
	domain.Update
	File FileMeta
}

func (c *PersonalChat) NewFileMessage(sender domain.UserID, file *FileMeta) (FileMessage, error) {
	if err := c.validateCanSend(sender); err != nil {
		return FileMessage{}, err
	}

	if err := validateFile(file); err != nil {
		return FileMessage{}, err
	}

	return FileMessage{
		Update: domain.Update{
			ChatID:   c.ChatID,
			SenderID: sender,
		},
		File: *file,
	}, nil
}

func (c *PersonalChat) DeleteFileMessage(sender domain.UserID, m *FileMessage, mode domain.DeleteMode) error {
	if err := c.validateCanSend(sender); err != nil {
		return err
	}

	if c.ChatID != m.ChatID {
		return domain.ErrUpdateNotFromChat
	}

	if m.DeletedFor(sender) {
		return domain.ErrUpdateDeleted
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
