package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type SecretGroupPhotoService struct {
	repo  repository.SecretGroupChatRepository
	files external.FileStorage
}

func NewSecretGroupPhotoService(repo repository.SecretGroupChatRepository) SecretGroupPhotoService {
	return SecretGroupPhotoService{
		repo: repo,
	}
}

func (s *SecretGroupPhotoService) UpdatePhoto(ctx context.Context, groupId, fileId uuid.UUID) (*dto.SecretGroupChatDTO, error) {
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

	if err := validatePhoto(file); err != nil {
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

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}

func (s *SecretGroupPhotoService) DeletePhoto(ctx context.Context, groupId uuid.UUID) (*dto.SecretGroupChatDTO, error) {
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

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}
