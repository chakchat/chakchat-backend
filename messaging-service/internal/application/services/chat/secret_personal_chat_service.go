package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
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

func (s *SecretPersonalChatService) CreateChat(ctx context.Context, userId, withUserId uuid.UUID) (*dto.SecretPersonalChatDTO, error) {
	domainMembers := [2]domain.UserID{domain.UserID(userId), domain.UserID(withUserId)}

	if err := s.validateChatNotExists(ctx, domainMembers); err != nil {
		return nil, err
	}

	chat, err := secpersonal.NewSecretPersonalChatService(domainMembers)

	if err != nil {
		if errors.Is(err, domain.ErrChatWithMyself) {
			return nil, ErrChatWithMyself
		}
		return nil, errors.Join(ErrInternal, err)
	}

	chat, err = s.repo.Create(ctx, chat)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	chatDto := dto.NewSecretPersonalChatDTO(chat)

	s.pub.PublishForUsers([]uuid.UUID{withUserId}, events.ChatCreated{
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
			return nil, ErrChatNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	chatDto := dto.NewSecretPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *SecretPersonalChatService) DeleteChat(ctx context.Context, chatId uuid.UUID) error {
	chat, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	// TODO: put other logic here after you decide what to do with messages
	// For now I think that messages with in deleted chat should be deleted by background task

	if err := s.repo.Delete(ctx, chat.ID); err != nil {
		return errors.Join(ErrInternal, err)
	}

	s.pub.PublishForUsers(dto.UUIDs(chat.Members[:]), events.ChatDeleted{
		ChatID: chatId,
	})

	return nil
}

func (s *SecretPersonalChatService) validateChatNotExists(ctx context.Context, members [2]domain.UserID) error {
	_, err := s.repo.FindByMembers(ctx, members)

	if err != nil && err != repository.ErrNotFound {
		return errors.Join(ErrInternal, err)
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return ErrChatAlreadyExists
	}

	return nil
}
