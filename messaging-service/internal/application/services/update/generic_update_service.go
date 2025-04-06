package update

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GenericUpdateService struct {
	txProvider storage.TxProvider
	chatRepo   repository.ChatterRepository
	updateRepo repository.GenericUpdateRepository
}

func NewGenericUpdateService(
	txProvider storage.TxProvider,
	chatRepo repository.ChatterRepository,
	updateRepo repository.GenericUpdateRepository,
) *GenericUpdateService {
	return &GenericUpdateService{
		txProvider: txProvider,
		chatRepo:   chatRepo,
		updateRepo: updateRepo,
	}
}

func (s *GenericUpdateService) GetUpdatesInRange(
	ctx context.Context, req request.GetUpdatesInRange,
) ([]services.GenericUpdate, error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	// It is such cringe, it should be refactored
	chat, err := s.chatRepo.FindChatter(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		return nil, err
	}
	if !chat.IsMember(domain.UserID(req.SenderID)) {
		return nil, domain.ErrUserNotMember
	}

	updates, err := s.updateRepo.GetInRange(
		ctx, tx,
		domain.UserID(req.SenderID),
		domain.ChatID(req.ChatID),
		domain.UpdateID(req.From),
		domain.UpdateID(req.To),
	)
	return updates, err
}

func (s *GenericUpdateService) GetUpdate(
	ctx context.Context, req request.GetUpdate,
) (*services.GenericUpdate, error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	// It is such cringe, it should be refactored
	chat, err := s.chatRepo.FindChatter(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		return nil, err
	}
	if !chat.IsMember(domain.UserID(req.SenderID)) {
		return nil, domain.ErrUserNotMember
	}

	updates, err := s.updateRepo.Get(
		ctx, tx,
		domain.UserID(req.SenderID),
		domain.ChatID(req.ChatID),
		domain.UpdateID(req.UpdateID),
	)
	return updates, err
}
