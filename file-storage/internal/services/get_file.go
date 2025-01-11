package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrFileNotFound = errors.New("file not found")

type FileMetaGetter interface {
	GetFileMeta(uuid.UUID) (*FileMeta, bool, error)
}

type GetFileService struct {
	getter FileMetaGetter
}

func NewGetFileService(getter FileMetaGetter) *GetFileService {
	return &GetFileService{
		getter: getter,
	}
}

func (s *GetFileService) GetFile(fileId uuid.UUID) (*FileMeta, error) {
	meta, ok, err := s.getter.GetFileMeta(fileId)
	if err != nil {
		return nil, fmt.Errorf("getting file metadata failed: %s", err)
	}
	if !ok {
		return nil, ErrFileNotFound
	}
	return meta, nil
}
