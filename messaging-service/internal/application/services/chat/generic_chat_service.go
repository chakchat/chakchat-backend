package chat

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type GenericChatService struct {
	repo repository.GenericChatRepository
}

func NewGenericChatService(repo repository.GenericChatRepository) *GenericChatService {
	return &GenericChatService{
		repo: repo,
	}
}

func (s *GenericChatService) GetByMemberID(ctx context.Context, memberID uuid.UUID) ([]services.GenericChat, error) {
	return s.repo.GetByMemberID(ctx, domain.UserID(memberID))
}

func (s *GenericChatService) GetByChatID(ctx context.Context, memberID uuid.UUID) (*services.GenericChat, error) {
	chat, err := s.repo.GetByChatID(ctx, domain.ChatID(memberID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, fmt.Errorf("getting generic chat failed: %s", err)
	}

	return chat, nil
}
