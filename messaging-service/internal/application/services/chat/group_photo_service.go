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

const (
	MaxGroupPhotoSize = 2 << 20
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
	pub   publish.Publisher
}

func NewGroupPhotoService(repo repository.GroupChatRepository,
	files external.FileStorage,
	pub publish.Publisher,
) *GroupPhotoService {
	return &GroupPhotoService{
		repo:  repo,
		files: files,
		pub:   pub,
	}
}

func (s *GroupPhotoService) UpdatePhoto(ctx context.Context, req request.UpdateGroupPhoto) (*dto.GroupChatDTO, error) {
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
		return nil, errors.Join(services.ErrInternal, err)
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.GroupInfoUpdated{
			ChatID:        gDto.ID,
			Name:          gDto.Name,
			Description:   gDto.Description,
			GroupPhotoURL: string(g.GroupPhoto),
		},
	)

	return &gDto, nil
}

func validatePhoto(photo *external.FileMeta) error {
	if photo.FileSize > MaxGroupPhotoSize {
		return services.ErrInvalidPhoto
	}

	if !groupPhotoMimes[photo.MimeType] {
		return services.ErrInvalidPhoto
	}

	return nil
}

func (s *GroupPhotoService) DeletePhoto(ctx context.Context, req request.DeleteGroupPhoto) (*dto.GroupChatDTO, error) {
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

	gDto := dto.NewGroupChatDTO(g)

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.GroupInfoUpdated{
			ChatID:        gDto.ID,
			Name:          gDto.Name,
			Description:   gDto.Description,
			GroupPhotoURL: string(g.GroupPhoto),
		},
	)

	return &gDto, nil
}
