package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/personal"
	"github.com/google/uuid"
)

type PersonalChatService struct {
	txProvider storage.TxProvider
	repo       repository.PersonalChatRepository
	pub        publish.Publisher
}

func NewPersonalChatService(
	txProvider storage.TxProvider, repo repository.PersonalChatRepository, pub publish.Publisher,
) *PersonalChatService {
	return &PersonalChatService{
		repo:       repo,
		pub:        pub,
		txProvider: txProvider,
	}
}

func (s *PersonalChatService) BlockChat(
	ctx context.Context, req request.BlockChat,
) (_ *dto.PersonalChatDTO, err error) {
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

	err = chat.BlockBy(domain.UserID(req.SenderID))

	if err != nil {
		return nil, err
	}

	if chat, err = s.repo.Update(ctx, tx, chat); err != nil {
		return nil, err
	}

	s.pub.PublishForUsers(
		services.GetReceivingMembers(chat.Members[:], domain.UserID(req.SenderID)),
		events.ChatBlocked{
			ChatID: req.ChatID,
		},
	)

	chatDto := dto.NewPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *PersonalChatService) UnblockChat(
	ctx context.Context, req request.UnblockChat,
) (_ *dto.PersonalChatDTO, err error) {
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

	err = chat.UnblockBy(domain.UserID(req.SenderID))

	if err != nil {
		return nil, err
	}

	if _, err := s.repo.Update(ctx, tx, chat); err != nil {
		return nil, err
	}

	s.pub.PublishForUsers(
		services.GetReceivingMembers(chat.Members[:], domain.UserID(req.SenderID)),
		events.ChatUnblocked{
			ChatID: req.ChatID,
		},
	)

	chatDto := dto.NewPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *PersonalChatService) CreateChat(
	ctx context.Context, req request.CreatePersonalChat,
) (_ *dto.PersonalChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	domainMembers := [2]domain.UserID{domain.UserID(req.SenderID), domain.UserID(req.MemberID)}

	if err := s.validateChatNotExists(ctx, domainMembers); err != nil {
		return nil, err
	}

	chat, err := personal.NewPersonalChat(domainMembers)

	if err != nil {
		return nil, err
	}

	chat, err = s.repo.Create(ctx, tx, chat)
	if err != nil {
		return nil, err
	}

	chatDto := dto.NewPersonalChatDTO(chat)

	s.pub.PublishForUsers([]uuid.UUID{req.MemberID}, events.ChatCreated{
		ChatID:   chatDto.ID,
		ChatType: events.ChatTypePersonal,
	})

	return &chatDto, nil
}

func (s *PersonalChatService) GetChatById(
	ctx context.Context, chatId uuid.UUID,
) (_ *dto.PersonalChatDTO, err error) {
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

	chatDto := dto.NewPersonalChatDTO(chat)
	return &chatDto, nil
}

func (s *PersonalChatService) DeleteChat(ctx context.Context, req request.DeleteChat) (err error) {
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

	s.pub.PublishForUsers(
		services.GetReceivingMembers(chat.Members[:], domain.UserID(req.SenderID)),
		events.ChatDeleted{
			ChatID: req.ChatID,
		},
	)

	return nil
}

func (s *PersonalChatService) validateChatNotExists(ctx context.Context, members [2]domain.UserID) (err error) {
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
