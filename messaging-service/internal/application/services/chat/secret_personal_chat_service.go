package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secpersonal"
	"github.com/google/uuid"
)

type SecretPersonalChatService struct {
	txProvider storage.TxProvider
	repo       repository.SecretPersonalChatRepository
	pub        publish.Publisher
}

func NewSecretPersonalChatService(
	txProvider storage.TxProvider,
	repo repository.SecretPersonalChatRepository,
	pub publish.Publisher,
) *SecretPersonalChatService {
	return &SecretPersonalChatService{
		repo:       repo,
		txProvider: txProvider,
		pub:        pub,
	}
}

func (s *SecretPersonalChatService) CreateChat(
	ctx context.Context, req request.CreateSecretPersonalChat,
) (_ *dto.SecretPersonalChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	domainMembers := [2]domain.UserID{domain.UserID(req.SenderID), domain.UserID(req.MemberID)}

	if err := s.validateChatNotExists(ctx, domainMembers); err != nil {
		return nil, err
	}

	chat, err := secpersonal.NewSecretPersonalChatService(domainMembers)

	if err != nil {
		return nil, err
	}

	chat, err = s.repo.Create(ctx, tx, chat)
	if err != nil {
		return nil, err
	}

	chatDto := dto.NewSecretPersonalChatDTO(chat)

	s.pub.PublishForReceivers(
		[]uuid.UUID{req.MemberID},
		events.TypeChatCreated,
		events.ChatCreated{
			SenderID: req.SenderID,
			Chat:     generic.FromSecretPersonalChatDTO(&chatDto),
		},
	)

	return &chatDto, nil
}

func (s *SecretPersonalChatService) GetChatById(
	ctx context.Context, chatId uuid.UUID,
) (_ *dto.SecretPersonalChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.repo.FindById(ctx, tx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	chatDto := dto.NewSecretPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *SecretPersonalChatService) SetExpiration(
	ctx context.Context, req request.SetExpiration,
) (_ *dto.SecretPersonalChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
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

	chat, err = s.repo.Update(ctx, tx, chat)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(chat.Members[:], domain.UserID(req.SenderID)),
		events.TypeChatExpirationSet,
		events.ExpirationSet{
			SenderID:   req.SenderID,
			ChatID:     req.ChatID,
			Expiration: req.Expiration,
		},
	)

	chatDto := dto.NewSecretPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *SecretPersonalChatService) DeleteChat(ctx context.Context, req request.DeleteChat) (err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
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

	if err := s.repo.Delete(ctx, tx, chat.ID); err != nil {
		return err
	}

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(chat.Members[:], domain.UserID(req.SenderID)),
		events.TypeChatDeleted,
		events.ChatDeleted{
			SenderID: req.SenderID,
			ChatID:   req.ChatID,
		},
	)

	return nil
}

func (s *SecretPersonalChatService) validateChatNotExists(
	ctx context.Context, members [2]domain.UserID,
) (err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return err
	}
	defer storage.FinishTx(ctx, tx, &err)

	_, err = s.repo.FindByMembers(ctx, tx, members)

	if err != nil && err != repository.ErrNotFound {
		return err
	}

	if !errors.Is(err, repository.ErrNotFound) {
		return services.ErrChatAlreadyExists
	}

	return nil
}
