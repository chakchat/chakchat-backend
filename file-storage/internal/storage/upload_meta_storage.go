package main

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat/backend/file-storage/internal/services"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadMeta struct {
	Id         uuid.UUID `gorm:"primaryKey"`
	Key        string
	FileName   string
	MimeType   string
	S3UploadId string
	FileId     uuid.UUID
}

type UploadMetaStorage struct {
	db *gorm.DB
}

func (s *UploadMetaStorage) Get(ctx context.Context, id uuid.UUID) (*services.UploadMeta, bool, error) {
	var meta UploadMeta
	if err := s.db.First(&meta, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &services.UploadMeta{
		PublicUploadId: meta.Id,
		Key:            meta.Key,
		FileName:       meta.FileName,
		MimeType:       meta.MimeType,
		S3UploadId:     meta.S3UploadId,
		FileId:         meta.FileId,
	}, true, nil
}

func (s *UploadMetaStorage) Store(ctx context.Context, m *services.UploadMeta) error {
	meta := UploadMeta{
		Id:         m.PublicUploadId,
		Key:        m.Key,
		FileName:   m.FileName,
		MimeType:   m.MimeType,
		S3UploadId: m.S3UploadId,
		FileId:     m.FileId,
	}

	if err := s.db.Create(meta).Error; err != nil {
		return err
	}
	return nil
}

func (s *UploadMetaStorage) Remove(ctx context.Context, id uuid.UUID) error {
	return s.db.Delete(&UploadMeta{}, id).Error
}
