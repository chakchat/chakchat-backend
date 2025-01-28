package services

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrAdminNotMember = errors.New("service: admin is not group member")

	ErrGroupNameEmpty   = errors.New("service: group name is empty")
	ErrGroupNameTooLong = errors.New("service: group name is too long")
	ErrGroupDescTooLong = errors.New("service: group description is too long")

	ErrUserAlreadyMember = errors.New("service: user is already a member of a chat")
	ErrMemberIsAdmin     = errors.New("service: group member is admin")

	ErrGroupPhotoEmpty = errors.New("service: group photo is empty")
)

type GroupChatService struct {
	repo repository.GroupChatRepository
}

func NewGroupChatService(repo repository.GroupChatRepository) *GroupChatService {
	return &GroupChatService{
		repo: repo,
	}
}

type CreateGroupRequest struct {
	Admin   uuid.UUID
	Members []uuid.UUID
	Name    string
}

func (s *GroupChatService) CreateGroup(ctx context.Context, req CreateGroupRequest) (*dto.GroupChatDTO, error) {
	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	group, err := domain.NewGroupChat(domain.UserID(req.Admin), members, req.Name)

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

	group, err = s.repo.Create(ctx, group)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	groupDto := dto.NewGroupChatDTO(group)
	return &groupDto, nil
}

type CreateSecretGroupRequest struct {
	Admin   uuid.UUID
	Members []uuid.UUID
	Name    string
}

func (s *GroupChatService) CreateSecretGroup(ctx context.Context, req CreateSecretGroupRequest) (*dto.GroupChatDTO, error) {
	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	group, err := domain.NewSecretGroupChat(domain.UserID(req.Admin), members, req.Name)

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

	group, err = s.repo.Create(ctx, group)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	groupDto := dto.NewGroupChatDTO(group)
	return &groupDto, nil
}

type UpdateGroupInfoRequest struct {
	ChatID      uuid.UUID
	Name        string
	Description string
}

func (s *GroupChatService) UpdateGroupInfo(ctx context.Context, req UpdateGroupInfoRequest) (*dto.GroupChatDTO, error) {
	group, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, errors.Join(ErrInternal, err)
	}

	err = group.UpdateInfo(req.Name, req.Description)

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

	group, err = s.repo.Update(ctx, group)
	if err != nil {
		return nil, errors.Join(ErrInternal, err)
	}

	groupDto := dto.NewGroupChatDTO(group)
	return &groupDto, nil
}
