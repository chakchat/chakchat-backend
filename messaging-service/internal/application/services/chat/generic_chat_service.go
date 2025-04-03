package chat

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type GenericChatService struct {
	txProvider storage.TxProvider
	repo       repository.GenericChatRepository
}

func NewGenericChatService(
	txProvider storage.TxProvider, repo repository.GenericChatRepository,
) *GenericChatService {
	return &GenericChatService{
		txProvider: txProvider,
		repo:       repo,
	}
}

func (s *GenericChatService) GetByMemberID(ctx context.Context, memberID uuid.UUID) (_ []services.GenericChat, err error) {
	tx, err := s.txProvider.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	return s.repo.GetByMemberID(ctx, tx, domain.UserID(memberID))
}

func (s *GenericChatService) GetByChatID(ctx context.Context, senderID, chatID uuid.UUID) (_ *services.GenericChat, err error) {
	tx, err := s.txProvider.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.repo.GetByChatID(ctx, tx, domain.ChatID(chatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, fmt.Errorf("getting generic chat failed: %s", err)
	}
	// TODO: Maybe it is cringe but refactor it later, for now I don't care
	if !slices.Contains(chat.Members, senderID) {
		return nil, domain.ErrUserNotMember
	}

	return chat, nil
}
