package external

import (
	"errors"
	"time"

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
	CreatedAt time.Time
}

type FileStorage interface {
	// Should return ErrFileNotFound if there are no such file
	GetById(uuid.UUID) (*FileMeta, error)
}
