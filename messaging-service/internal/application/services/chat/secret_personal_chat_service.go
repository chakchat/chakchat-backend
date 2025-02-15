package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secpersonal"
	"github.com/google/uuid"
)

type SecretPersonalChatService struct {
	repo repository.SecretPersonalChatRepository
	pub  publish.Publisher
}

func NewSecretPersonalChatService(repo repository.SecretPersonalChatRepository,
	pub publish.Publisher,
) *SecretPersonalChatService {
	return &SecretPersonalChatService{
		repo: repo,
	}
}

func (s *SecretPersonalChatService) CreateChat(ctx context.Context, req request.CreateSecretPersonalChat) (*dto.SecretPersonalChatDTO, error) {
	domainMembers := [2]domain.UserID{domain.UserID(req.SenderID), domain.UserID(req.MemberID)}

	if err := s.validateChatNotExists(ctx, domainMembers); err != nil {
		return nil, err
	}

	chat, err := secpersonal.NewSecretPersonalChatService(domainMembers)

	if err != nil {
		return nil, err
	}

	chat, err = s.repo.Create(ctx, chat)
	if err != nil {
		return nil, err
	}

	chatDto := dto.NewSecretPersonalChatDTO(chat)

	s.pub.PublishForUsers([]uuid.UUID{req.MemberID}, events.ChatCreated{
		ChatID:   chatDto.ID,
		ChatType: events.ChatTypePersonal,
	})

	return &chatDto, nil
}

func (s *SecretPersonalChatService) GetChatById(ctx context.Context, chatId uuid.UUID,
) (*dto.SecretPersonalChatDTO, error) {
	chat, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	chatDto := dto.NewSecretPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *SecretPersonalChatService) SetExpiration(ctx context.Context, req request.SetExpiration) (*dto.SecretPersonalChatDTO, error) {
	chat, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	err = chat.SetExpiration(domain.UserID(req.SenderID), req.Expiration)
	if err != nil {
		return nil, err
	}

	chat, err = s.repo.Update(ctx, chat)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForUsers(
		services.GetReceivingMembers(chat.Members[:], domain.UserID(req.SenderID)),
		events.ExpirationSet{
			ChatID:     req.ChatID,
			SenderID:   req.SenderID,
			Expiration: req.Expiration,
		},
	)

	chatDto := dto.NewSecretPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *SecretPersonalChatService) DeleteChat(ctx context.Context, req request.DeleteChat) error {
	chat, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return services.ErrChatNotFound
		}
		return err
	}

	err = chat.Delete(domain.UserID(req.SenderID))
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, chat.ID); err != nil {
		return err
	}

	s.pub.PublishForUsers(
		services.GetReceivingMembers(chat.Members[:], domain.UserID(req.SenderID)),
		events.ChatDeleted{
			ChatID: req.ChatID,
		},
	)

	return nil
}

func (s *SecretPersonalChatService) validateChatNotExists(ctx context.Context, members [2]domain.UserID) error {
	_, err := s.repo.FindByMembers(ctx, members)

	if err != nil && err != repository.ErrNotFound {
		return err
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return services.ErrChatAlreadyExists
	}

	return nil
}
