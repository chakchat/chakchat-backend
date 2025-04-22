package chat

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GenericChatService struct {
	txProvider storage.TxProvider
	chatRepo   repository.GenericChatRepository
	updaterepo repository.GenericUpdateRepository
}

func NewGenericChatService(
	txProvider storage.TxProvider,
	chatRepo repository.GenericChatRepository,
	updateRepo repository.GenericUpdateRepository,
) *GenericChatService {
	return &GenericChatService{
		txProvider: txProvider,
		chatRepo:   chatRepo,
		updaterepo: updateRepo,
	}
}

func (s *GenericChatService) GetByMemberID(ctx context.Context, memberID uuid.UUID, opts ...request.GetChatOption) (_ []generic.Chat, err error) {
	opt := request.NewGetChatOptions(opts...)

	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chats, err := s.chatRepo.GetByMemberID(ctx, tx, domain.UserID(memberID))
	if err != nil {
		return nil, err
	}

	if opt.LoadLastUpdateID {
		for _, chat := range chats {
			if err = s.fillLastUpdateID(ctx, tx, &chat); err != nil {
				return nil, err
			}
		}
	}

	if opt.LoadPreviewCount > 0 {
		for i := range chats {
			if err = s.fillPreview(ctx, tx, &chats[i], memberID, opt.LoadPreviewCount); err != nil {
				return nil, err
			}
		}
	}

	return chats, nil
}

func (s *GenericChatService) GetByChatID(ctx context.Context, senderID, chatID uuid.UUID, opts ...request.GetChatOption) (_ *generic.Chat, err error) {
	opt := request.NewGetChatOptions(opts...)

	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.chatRepo.GetByChatID(ctx, tx, domain.ChatID(chatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, fmt.Errorf("getting generic chat failed: %s", err)
	}

	if !slices.Contains(chat.Members, senderID) {
		return nil, domain.ErrUserNotMember
	}

	if opt.LoadLastUpdateID {
		if err = s.fillLastUpdateID(ctx, tx, chat); err != nil {
			return nil, err
		}
	}

	if opt.LoadPreviewCount > 0 {
		if err = s.fillPreview(ctx, tx, chat, senderID, opt.LoadPreviewCount); err != nil {
			return nil, err
		}
	}

	return chat, nil
}

func (s *GenericChatService) fillLastUpdateID(ctx context.Context, tx pgx.Tx, chat *generic.Chat) error {
	lastUpdateID, err := s.updaterepo.GetLastUpdateID(ctx, tx, domain.ChatID(chat.ChatID))
	if err != nil {
		return fmt.Errorf("fill last UpdateID: %w", err)
	}

	cp := int64(lastUpdateID)
	chat.LastUpdateID = &cp
	return nil
}

func (s *GenericChatService) fillPreview(
	ctx context.Context, tx pgx.Tx, chat *generic.Chat, senderID uuid.UUID, previewCount int,
) error {
	updates, err := s.updaterepo.FetchLast(
		ctx, tx, domain.UserID(senderID), domain.ChatID(chat.ChatID),
		repository.WithFetchLastCount(previewCount),
		repository.WithFetchLastOptions(repository.FetchLastModeMessages),
	)
	if err != nil {
		return err
	}

	chat.UpdatePreview = updates
	return nil
}
