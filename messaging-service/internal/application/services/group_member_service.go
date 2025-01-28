package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type GroupMemberService struct {
	repo repository.GroupChatRepository
}

func NewGroupMembersService(repo repository.GroupChatRepository) *GroupMemberService {
	return &GroupMemberService{
		repo: repo,
	}
}

func (s *GroupMemberService) AddMember(ctx context.Context, chatId, userId uuid.UUID) (*dto.GroupChatDTO, error) {
	group, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
	}

	err = group.AddMember(domain.UserID(userId))

	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyMember) {
			return nil, ErrUserAlreadyMember
		}
		return nil, errors.Join(ErrInternal, err)
	}

	group, err = s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	groupDto := dto.NewGroupChatDTO(group)
	return &groupDto, nil
}

func (s *GroupMemberService) DeleteMember(ctx context.Context, chatId, memberId uuid.UUID) (*dto.GroupChatDTO, error) {
	group, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
	}

	err = group.DeleteMember(domain.UserID(memberId))

	switch {
	case errors.Is(err, domain.ErrMemberIsAdmin):
		return nil, ErrMemberIsAdmin
	case errors.Is(err, domain.ErrUserNotMember):
		return nil, ErrUserNotMember
	case err != nil:
		return nil, errors.Join(ErrInternal, err)
	}

	group, err = s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	groupDto := dto.NewGroupChatDTO(group)
	return &groupDto, nil
}
