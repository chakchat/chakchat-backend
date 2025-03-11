package updatex

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/update"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type SecretUpdateService struct {
	genChatRepo     repository.GenericChatRepository
	secPersonalSrvc *update.SecretPersonalUpdateService
	secGroupSrvc    *update.SecretGroupUpdateService
}

func NewSecretUpdateService(
	genChatRepo repository.GenericChatRepository,
	personalSrvc *update.SecretPersonalUpdateService,
	groupSrvc *update.SecretGroupUpdateService,
) *SecretUpdateService {
	return &SecretUpdateService{
		genChatRepo:     genChatRepo,
		secPersonalSrvc: personalSrvc,
		secGroupSrvc:    groupSrvc,
	}
}

func (s *SecretUpdateService) SendSecretUpdate(ctx context.Context, req request.SendSecretUpdate) (*dto.SecretUpdateDTO, error) {
	chatType, err := s.getChatType(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypeSecretPersonal:
		return s.secPersonalSrvc.SendSecretUpdate(ctx, req)
	case services.ChatTypeSecretGroup:
		return s.secGroupSrvc.SendSecretUpdate(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot send secret update in the chat of type: %s", chatType))
	}
}

func (s *SecretUpdateService) DeleteSecretUpdate(ctx context.Context, req request.DeleteSecretUpdate) (*dto.UpdateDeletedDTO, error) {
	chatType, err := s.getChatType(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypeSecretPersonal:
		return s.secPersonalSrvc.DeleteSecretUpdate(ctx, req)
	case services.ChatTypeSecretGroup:
		return s.secGroupSrvc.DeleteSecretUpdate(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot delete secret update in the chat of type: %s", chatType))
	}
}

func (s *SecretUpdateService) getChatType(ctx context.Context, id uuid.UUID) (string, error) {
	chatType, err := s.genChatRepo.GetChatType(ctx, domain.ChatID(id))
	if errors.Is(err, repository.ErrNotFound) {
		return "", services.ErrChatNotFound
	}
	return chatType, err
}
