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
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type SecretGroupPhotoService struct {
	txProvider storage.TxProvider
	repo       repository.SecretGroupChatRepository
	files      external.FileStorage
	pub        publish.Publisher
}

func NewSecretGroupPhotoService(
	txProvider storage.TxProvider,
	repo repository.SecretGroupChatRepository,
	files external.FileStorage,
	pub publish.Publisher,
) *SecretGroupPhotoService {
	return &SecretGroupPhotoService{
		repo:       repo,
		files:      files,
		pub:        pub,
		txProvider: txProvider,
	}
}

func (s *SecretGroupPhotoService) UpdatePhoto(
	ctx context.Context, req request.UpdateGroupPhoto,
) (_ *dto.SecretGroupChatDTO, err error) {
	tx, err := s.txProvider.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	g, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	file, err := s.files.GetById(ctx, req.FileID)
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, services.ErrFileNotFound
		}
		return nil, err
	}

	if err := validatePhoto(file); err != nil {
		return nil, err
	}

	err = g.UpdatePhoto(domain.UserID(req.SenderID), domain.URL(file.FileUrl))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

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

func (s *SecretGroupPhotoService) DeletePhoto(
	ctx context.Context, req request.DeleteGroupPhoto,
) (_ *dto.SecretGroupChatDTO, err error) {
	tx, err := s.txProvider.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	g, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, services.ErrFileNotFound
		}
		return nil, err
	}

	err = g.DeletePhoto(domain.UserID(req.SenderID))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

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
