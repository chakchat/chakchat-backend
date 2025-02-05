package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/query"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/group"
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
	pub  publish.Publisher
}

func NewGroupChatService(repo repository.GroupChatRepository, pub publish.Publisher) *GroupChatService {
	return &GroupChatService{
		repo: repo,
		pub:  pub,
	}
}

func (s *GroupChatService) CreateGroup(ctx context.Context, req query.CreateGroupRequest) (*dto.GroupChatDTO, error) {
	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	g, err := group.NewGroupChat(domain.UserID(req.Admin), members, req.Name)

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

	gDto := dto.NewGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.ChatCreated{
		ChatID:   uuid.UUID(gDto.ID),
		ChatType: events.ChatTypeGroup,
	})

	return &gDto, nil
}

func (s *GroupChatService) UpdateGroupInfo(ctx context.Context, req query.UpdateGroupInfoRequest) (*dto.GroupChatDTO, error) {
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

	gDto := dto.NewGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.GroupInfoUpdated{
		ChatID:        gDto.ID,
		Name:          gDto.Name,
		Description:   gDto.Description,
		GroupPhotoURL: string(g.GroupPhoto),
	})

	return &gDto, nil
}

func (s *GroupChatService) DeleteGroup(ctx context.Context, chatId uuid.UUID) error {
	g, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		return errors.Join(ErrInternal, err)
	}

	// TODO: put other logic here after you decide what to do with messages

	if err := s.repo.Delete(ctx, g.ID); err != nil {
		return errors.Join(ErrInternal, err)
	}

	s.pub.PublishForUsers(dto.UUIDs(g.Members), events.ChatDeleted{
		ChatID: chatId,
	})

	return nil
}
