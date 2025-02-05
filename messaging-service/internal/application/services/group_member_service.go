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
	g, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
	}

	err = g.AddMember(domain.UserID(userId))

	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyMember) {
			return nil, ErrUserAlreadyMember
		}
		return nil, errors.Join(ErrInternal, err)
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)
	return &gDto, nil
}

func (s *GroupMemberService) DeleteMember(ctx context.Context, chatId, memberId uuid.UUID) (*dto.GroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
	}

	err = g.DeleteMember(domain.UserID(memberId))

	switch {
	case errors.Is(err, domain.ErrMemberIsAdmin):
		return nil, ErrMemberIsAdmin
	case errors.Is(err, domain.ErrUserNotMember):
		return nil, ErrUserNotMember
	case err != nil:
		return nil, errors.Join(ErrInternal, err)
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)
	return &gDto, nil
}
