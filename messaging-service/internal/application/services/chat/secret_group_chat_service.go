package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secgroup"
	"github.com/google/uuid"
)

type SecretGroupChatService struct {
	repo repository.SecretGroupChatRepository
}

func NewSecretGroupChatService(repo repository.SecretGroupChatRepository) *SecretGroupChatService {
	return &SecretGroupChatService{
		repo: repo,
	}
}

type CreateSecretGroupRequest struct {
	Admin   uuid.UUID
	Members []uuid.UUID
	Name    string
}

func (s *SecretGroupChatService) CreateGroup(ctx context.Context, req CreateGroupRequest) (*dto.SecretGroupChatDTO, error) {
	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	g, err := secgroup.NewSecretGroupChat(domain.UserID(req.Admin), members, req.Name)

	switch {
	case errors.Is(err, domain.ErrAdminNotMember):
		return nil, ErrAdminNotMember
	case errors.Is(err, domain.ErrGroupNameEmpty):
		return nil, ErrGroupNameEmpty
	case errors.Is(err, domain.ErrGroupNameTooLong):
		return nil, ErrGroupNameTooLong
	case err != nil:
		return nil, errors.Join(ErrInternal, err)
	}

	g, err = s.repo.Create(ctx, g)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}

type UpdateSecretGroupInfoRequest struct {
	ChatID      uuid.UUID
	Name        string
	Description string
}

func (s *SecretGroupChatService) UpdateGroupInfo(ctx context.Context, req UpdateGroupInfoRequest) (*dto.SecretGroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	err = g.UpdateInfo(req.Name, req.Description)

	switch {
	case errors.Is(err, domain.ErrGroupNameEmpty):
		return nil, ErrGroupNameEmpty
	case errors.Is(err, domain.ErrGroupNameTooLong):
		return nil, ErrGroupNameTooLong
	case errors.Is(err, domain.ErrGroupDescTooLong):
		return nil, ErrGroupDescTooLong
	case err != nil:
		return nil, errors.Join(ErrInternal, err)
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}

func (s *SecretGroupChatService) DeleteGroup(ctx context.Context, chatId uuid.UUID) error {
	chat, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	// TODO: put other logic here after you decide what to do with messages

	if err := s.repo.Delete(ctx, chat.ID); err != nil {
		return errors.Join(ErrInternal, err)
	}
	return nil
}
