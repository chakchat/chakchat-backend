package update

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
)

type SecretPersonalUpdateService struct {
	txProvider       storage.TxProvider
	chatRepo         repository.SecretPersonalChatRepository
	secretUpdateRepo repository.SecretUpdateRepository
	pub              publish.Publisher
}

func NewSecretPersonalUpdateService(
	txProvider storage.TxProvider,
	chatRepo repository.SecretPersonalChatRepository,
	secretUpdateRepo repository.SecretUpdateRepository,
	pub publish.Publisher,
) *SecretPersonalUpdateService {
	return &SecretPersonalUpdateService{
		txProvider:       txProvider,
		chatRepo:         chatRepo,
		secretUpdateRepo: secretUpdateRepo,
		pub:              pub,
	}
}

func (s *SecretPersonalUpdateService) SendSecretUpdate(
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
	updateDto := dto.NewSecretUpdateDTO(update)

	err = s.pub.PublishForReceivers(
		ctx,
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &update.Update),
		events.TypeUpdate,
		generic.FromSecretUpdateDTO(&updateDto),
	)
	if err != nil {
		return nil, err
	}

	return &updateDto, nil
}

func (s *SecretPersonalUpdateService) DeleteSecretUpdate(
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
	deletedDto := dto.NewUpdateDeletedDTO(deleted)

	err = s.pub.PublishForReceivers(
		ctx,
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &update.Update),
		events.TypeUpdate,
		generic.FromUpdateDeletedDTO(&deletedDto),
	)
	if err != nil {
		return nil, err
	}

	return &deletedDto, nil
}
