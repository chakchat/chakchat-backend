package update

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
	"github.com/google/uuid"
)

type SecretPersonalUpdateService struct {
	chatRepo         repository.SecretPersonalChatRepository
	secretUpdateRepo repository.SecretUpdateRepository
	pub              publish.Publisher
}

func NewSecretPersonalUpdateService(
	chatRepo repository.SecretPersonalChatRepository,
	secretUpdateRepo repository.SecretUpdateRepository,
	pub publish.Publisher,
) *SecretPersonalUpdateService {
	return &SecretPersonalUpdateService{
		chatRepo:         chatRepo,
		secretUpdateRepo: secretUpdateRepo,
		pub:              pub,
	}
}

func (s *SecretPersonalUpdateService) SendSecretUpdate(ctx context.Context, req request.SendSecretUpdate) (*dto.SecretUpdateDTO, error) {
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
