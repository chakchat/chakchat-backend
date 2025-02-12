package domain

import (
	"github.com/google/uuid"
)

const (
	MaxFileSize = 1 << 30
)

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

func NewFileMessage(chat Chatter, sender UserID, file *FileMeta, replyTo *Message) (*FileMessage, error) {
	if err := chat.ValidateCanSend(sender); err != nil {
		return nil, err
	}

	if err := validateFile(file); err != nil {
		return nil, err
	}

	if replyTo != nil {
		if err := validateCanReply(chat, sender, replyTo); err != nil {
			return nil, err
		}
	}

	return &FileMessage{
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

func (m *FileMessage) Forward(chat Chatter, sender UserID, destChat Chatter) (*FileMessage, error) {
	if !chat.IsMember(sender) {
		return nil, ErrUserNotMember
	}
	if err := destChat.ValidateCanSend(sender); err != nil {
		return nil, err
	}

	return &FileMessage{
		Message: Message{
			Update: Update{
				ChatID:   destChat.ChatID(),
				SenderID: sender,
			},
			Forwarded: true,
		},
		File: m.File,
	}, nil
}

func validateFile(file *FileMeta) error {
	if file.FileSize > MaxFileSize {
		return ErrFileTooBig
	}
	return nil
}
