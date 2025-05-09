package external

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type FileMeta struct {
	FileId    uuid.UUID
	FileName  string
	MimeType  string
	FileSize  int64
	FileUrl   string
	CreatedAt int64
}

type FileStorage interface {
	// Should return ErrFileNotFound if there are no such file
	GetById(context.Context, uuid.UUID) (*FileMeta, error)
}
