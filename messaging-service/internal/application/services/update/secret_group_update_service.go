package update

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
	"github.com/google/uuid"
)

type SecretGroupUpdateService struct {
	txProvider       storage.TxProvider
	chatRepo         repository.SecretGroupChatRepository
	secretUpdateRepo repository.SecretUpdateRepository
	pub              publish.Publisher
}

func NewSecretGroupUpdateService(
	txProvider storage.TxProvider,
	chatRepo repository.SecretGroupChatRepository,
	secretUpdateRepo repository.SecretUpdateRepository,
	pub publish.Publisher,
) *SecretGroupUpdateService {
	return &SecretGroupUpdateService{
		chatRepo:         chatRepo,
		secretUpdateRepo: secretUpdateRepo,
		txProvider:       txProvider,
		pub:              pub,
	}
}

func (s *SecretGroupUpdateService) SendSecretUpdate(
	ctx context.Context, req request.SendSecretUpdate,
) (_ *dto.SecretUpdateDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.chatRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	update, err := domain.NewSecretUpdate(chat, domain.UserID(req.SenderID), domain.SecretData{
		KeyHash: domain.SecretKeyHash(req.KeyHash),
		Payload: req.Payload,
		IV:      req.InitializationVector,
	})
	if err != nil {
		return nil, err
	}

	update, err = s.secretUpdateRepo.CreateSecretUpdate(ctx, tx, update)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &update.Update),
		events.SecretUpdateSent{
			ChatID:               uuid.UUID(update.ChatID),
			UpdateID:             int64(update.UpdateID),
			SenderID:             uuid.UUID(update.SenderID),
			Payload:              update.Data.Payload,
			InitializationVector: update.Data.Payload,
			KeyHash:              string(update.Data.KeyHash),
			CreatedAt:            int64(update.CreatedAt),
		},
	)

	updateDto := dto.NewSecretUpdateDTO(update)
	return &updateDto, nil
}

func (s *SecretGroupUpdateService) DeleteSecretUpdate(
	ctx context.Context, req request.DeleteSecretUpdate,
) (_ *dto.UpdateDeletedDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.chatRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	update, err := s.secretUpdateRepo.FindSecretUpdate(
		ctx, tx, domain.ChatID(req.ChatID), domain.UpdateID(req.SecretUpdateID),
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrSecretUpdateNotFound
		}
		return nil, err
	}

	err = update.Delete(chat, domain.UserID(req.SenderID))
	if err != nil {
		return nil, err
	}

	update.Deleted[len(update.Deleted)-1], err = s.secretUpdateRepo.CreateUpdateDeleted(
		ctx, tx, update.Deleted[len(update.Deleted)-1],
	)
	if err != nil {
		return nil, err
	}

	// I do not delete updates because it may cause incosistency.
	// Add triggers and change on delete behavior on foreighn keys before deleting physically

	deleted := update.Deleted[len(update.Deleted)-1]
	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &update.Update),
		events.UpdateDeleted{
			ChatID:     uuid.UUID(deleted.ChatID),
			UpdateID:   int64(deleted.UpdateID),
			SenderID:   uuid.UUID(deleted.SenderID),
			DeletedID:  int64(deleted.DeletedID),
			DeleteMode: string(deleted.Mode),
			CreatedAt:  int64(deleted.CreatedAt),
		},
	)

	deletedDto := dto.NewUpdateDeletedDTO(deleted)
	return &deletedDto, nil
}
