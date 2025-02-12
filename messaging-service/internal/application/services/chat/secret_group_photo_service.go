package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type SecretGroupPhotoService struct {
	repo  repository.SecretGroupChatRepository
	files external.FileStorage
	pub   publish.Publisher
}

func NewSecretGroupPhotoService(repo repository.SecretGroupChatRepository,
	files external.FileStorage,
	pub publish.Publisher,
) SecretGroupPhotoService {
	return SecretGroupPhotoService{
		repo:  repo,
		files: files,
		pub:   pub,
	}
}

func (s *SecretGroupPhotoService) UpdatePhoto(ctx context.Context, req request.UpdateGroupPhoto) (*dto.SecretGroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	file, err := s.files.GetById(req.FileID)
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, services.ErrFileNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	if err := validatePhoto(file); err != nil {
		return nil, err
	}

	err = g.UpdatePhoto(domain.UserID(req.SenderID), domain.URL(file.FileUrl))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.GroupInfoUpdated{
		ChatID:        gDto.ID,
		Name:          gDto.Name,
		Description:   gDto.Description,
		GroupPhotoURL: string(g.GroupPhoto),
	})

	return &gDto, nil
}

func (s *SecretGroupPhotoService) DeletePhoto(ctx context.Context, req request.DeleteGroupPhoto) (*dto.SecretGroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, services.ErrFileNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	err = g.DeletePhoto(domain.UserID(req.SenderID))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.GroupInfoUpdated{
		ChatID:        gDto.ID,
		Name:          gDto.Name,
		Description:   gDto.Description,
		GroupPhotoURL: string(g.GroupPhoto),
	})

	return &gDto, nil
}
