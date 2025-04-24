package dto

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type FileMetaDTO struct {
	FileId    uuid.UUID
	FileName  string
	MimeType  string
	FileSize  int64
	FileURL   string
	CreatedAt int64
}

func NewFileMetaDTO(f *domain.FileMeta) FileMetaDTO {
	return FileMetaDTO{
		FileId:    f.FileId,
		FileName:  f.FileName,
		MimeType:  f.MimeType,
		FileSize:  f.FileSize,
		FileURL:   string(f.FileURL),
		CreatedAt: int64(f.CreatedAt),
	}
}

type FileMessageDTO struct {
	ChatID   uuid.UUID
	UpdateID int64
	SenderID uuid.UUID

	File    FileMetaDTO
	ReplyTo *int64

	CreatedAt int64
}

func NewFileMessageDTO(m *domain.FileMessage) FileMessageDTO {
	var replyTo *int64
	if m.ReplyTo != nil {
		cp := int64(*m.ReplyTo)
		replyTo = &cp
	}

	return FileMessageDTO{
		ChatID:    uuid.UUID(m.ChatID),
		UpdateID:  int64(m.UpdateID),
		SenderID:  uuid.UUID(m.SenderID),
		File:      NewFileMetaDTO(&m.File),
		ReplyTo:   replyTo,
		CreatedAt: int64(m.CreatedAt),
	}
}
