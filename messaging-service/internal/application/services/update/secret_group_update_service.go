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
	chatRepo         repository.SecretGroupChatRepository
	secretUpdateRepo repository.SecretUpdateRepository
	txProvider       storage.TxProvider
	pub              publish.Publisher
}

func NewSecretGroupUpdateService(
	chatRepo repository.SecretGroupChatRepository,
	secretUpdateRepo repository.SecretUpdateRepository,
	txProvider storage.TxProvider,
	pub publish.Publisher,
) *SecretGroupUpdateService {
	return &SecretGroupUpdateService{
		chatRepo:         chatRepo,
		secretUpdateRepo: secretUpdateRepo,
		txProvider:       txProvider,
		pub:              pub,
	}
}

func (s *SecretGroupUpdateService) SendSecretUpdate(ctx context.Context, req request.SendSecretUpdate) (*dto.SecretUpdateDTO, error) {
	chat, err := s.chatRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	update, err := domain.NewSecretUpdate(chat, domain.UserID(req.SenderID), domain.SecretData{
		KeyID:   domain.SecretKeyID(req.KeyID),
		Payload: req.Payload,
		IV:      req.InitializationVector,
	})
	if err != nil {
		return nil, err
	}

	update, err = s.secretUpdateRepo.CreateSecretUpdate(ctx, update)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &update.Update),
		events.SecretUpdateSent{
			ChatID:               uuid.UUID(update.ChatID),
			UpdateID:             int64(update.UpdateID),
			SenderID:             uuid.UUID(update.SenderID),
			Payload:              update.Data.Payload,
			InitializationVector: update.Data.Payload,
			KeyID:                uuid.UUID(update.Data.KeyID),
			CreatedAt:            int64(update.CreatedAt),
		},
	)

	updateDto := dto.NewSecretUpdateDTO(update)
	return &updateDto, nil
}

func (s *SecretGroupUpdateService) DeleteSecretUpdate(ctx context.Context, req request.DeleteSecretUpdate) (*dto.UpdateDeletedDTO, error) {
	chat, err := s.chatRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	update, err := s.secretUpdateRepo.FindSecretUpdate(ctx, domain.ChatID(req.ChatID), domain.UpdateID(req.SecretUpdateID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrSecretUpdateNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	err = update.Delete(chat, domain.UserID(req.SenderID))
	if err != nil {
		return nil, err
	}

	err = storage.RunInTx(ctx, s.txProvider, func(ctx context.Context) error {
		update.Deleted[len(update.Deleted)-1], err = s.secretUpdateRepo.CreateUpdateDeleted(ctx, update.Deleted[len(update.Deleted)-1])
		if err != nil {
			return err
		}

		if update.DeletedForAll() {
			err = s.secretUpdateRepo.DeleteSecretUpdate(ctx, update.ChatID, update.UpdateID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

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
