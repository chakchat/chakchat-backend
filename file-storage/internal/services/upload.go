package services

import (
	"io"

	"github.com/google/uuid"
)

type UploadFileRequest struct {
	FileName string
	MimeType string
	FileSize int64

	File io.Reader
}

type FileMeta struct {
	FileName string
	MimeType string
	FileSize int64
	FileId   uuid.UUID
	FileUrl  string
}
