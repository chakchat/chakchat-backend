package storage

import (
	"context"
	"errors"
	"time"

	"github.com/chakchat/chakchat/backend/file-storage/internal/services"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileMeta struct {
	FileId    uuid.UUID `gorm:"primaryKey"`
	FileName  string
	MimeType  string
	FileSize  int64
	FileUrl   string
	CreatedAt time.Time
}

type FileMetaStorage struct {
	db *gorm.DB
}

func NewFileMetaStorage(db *gorm.DB) *FileMetaStorage {
	return &FileMetaStorage{
		db: db,
	}
}

func (s *FileMetaStorage) GetFileMeta(ctx context.Context, id uuid.UUID) (*services.FileMeta, bool, error) {
	var meta FileMeta
	if err := s.db.First(&meta, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &services.FileMeta{
		FileName:  meta.FileName,
		MimeType:  meta.MimeType,
		FileSize:  meta.FileSize,
		FileId:    meta.FileId,
		FileUrl:   meta.FileUrl,
		CreatedAt: meta.CreatedAt,
	}, true, nil
}

func (s *FileMetaStorage) Store(ctx context.Context, m *services.FileMeta) error {
	meta := FileMeta{
		FileId:    m.FileId,
		FileName:  m.FileName,
		MimeType:  m.MimeType,
		FileSize:  m.FileSize,
		FileUrl:   m.FileUrl,
		CreatedAt: m.CreatedAt,
	}

	if err := s.db.Create(meta).Error; err != nil {
		return err
	}
	return nil
}
