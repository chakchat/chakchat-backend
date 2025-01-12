package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrFileNotFound = errors.New("file not found")

type FileMetaGetter interface {
	GetFileMeta(context.Context, uuid.UUID) (*FileMeta, bool, error)
}

type GetFileService struct {
	getter FileMetaGetter
}

func NewGetFileService(getter FileMetaGetter) *GetFileService {
	return &GetFileService{
		getter: getter,
	}
}

func (s *GetFileService) GetFile(ctx context.Context, fileId uuid.UUID) (*FileMeta, error) {
	meta, ok, err := s.getter.GetFileMeta(ctx, fileId)
	if err != nil {
		return nil, fmt.Errorf("getting file metadata failed: %s", err)
	}
	if !ok {
		return nil, ErrFileNotFound
	}
	return meta, nil
}
