package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type SecretGroupMemberService struct {
	repo repository.SecretGroupChatRepository
	pub  publish.Publisher
}

func NewSecretGroupMemberService(repo repository.SecretGroupChatRepository, pub publish.Publisher) *SecretGroupMemberService {
	return &SecretGroupMemberService{
		repo: repo,
		pub:  pub,
	}
}

func (s *SecretGroupMemberService) AddMember(ctx context.Context, chatId, userId uuid.UUID) (*dto.SecretGroupChatDTO, error) {
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

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.GroupMemberAdded{
		ChatID:   chatId,
		MemberID: userId,
	})

	return &gDto, nil
}

func (s *SecretGroupMemberService) DeleteMember(ctx context.Context, chatId, memberId uuid.UUID) (*dto.SecretGroupChatDTO, error) {
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

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.GroupMemberAdded{
		ChatID:   chatId,
		MemberID: memberId,
	})

	return &gDto, nil
}
