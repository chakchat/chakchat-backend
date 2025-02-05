package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

const (
	MaxGroupPhotoSize = 2 << 20
)

var (
	ErrFileNotFound = errors.New("service: file not found")
	ErrInvalidPhoto = errors.New("service: invalid photo")
)

var groupPhotoMimes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
	"image/heif": true,
	"image/heic": true,
}

type GroupPhotoService struct {
	repo  repository.GroupChatRepository
	files external.FileStorage
}

func NewGroupPhotoService(repo repository.GroupChatRepository, files external.FileStorage) *GroupPhotoService {
	return &GroupPhotoService{
		repo:  repo,
		files: files,
	}
}

func (s *GroupPhotoService) UpdatePhoto(ctx context.Context, groupId, fileId uuid.UUID) (*dto.GroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(groupId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	file, err := s.files.GetById(fileId)
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	if err := s.validatePhoto(file); err != nil {
		return nil, err
	}

	err = g.UpdatePhoto(domain.URL(file.FileUrl))

	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)
	return &gDto, nil
}

func (s *GroupPhotoService) validatePhoto(photo *external.FileMeta) error {
	if photo.FileSize > MaxGroupPhotoSize {
		return ErrInvalidPhoto
	}

	if !groupPhotoMimes[photo.MimeType] {
		return ErrInvalidPhoto
	}

	return nil
}

func (s *GroupPhotoService) DeletePhoto(ctx context.Context, groupId uuid.UUID) (*dto.GroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(groupId))
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	err = g.DeletePhoto()

	if err != nil {
		if errors.Is(err, domain.ErrGroupPhotoEmpty) {
			return nil, ErrGroupPhotoEmpty
		}
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)
	return &gDto, nil
}
