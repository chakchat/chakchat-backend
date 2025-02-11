package chat

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
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/group"
	"github.com/google/uuid"
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

func (s *GroupChatService) CreateGroup(ctx context.Context, req request.CreateGroup) (*dto.GroupChatDTO, error) {
	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	g, err := group.NewGroupChat(domain.UserID(req.Admin), members, req.Name)

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Create(ctx, g)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.ChatCreated{
		ChatID:   uuid.UUID(gDto.ID),
		ChatType: events.ChatTypeGroup,
	})

	return &gDto, nil
}

func (s *GroupChatService) UpdateGroupInfo(ctx context.Context, req request.UpdateGroupInfo) (*dto.GroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	err = g.UpdateInfo(req.Name, req.Description)

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
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
			return services.ErrChatNotFound
		}
		return errors.Join(services.ErrInternal, err)
	}

	// TODO: put other logic here after you decide what to do with messages

	if err := s.repo.Delete(ctx, g.ID); err != nil {
		return errors.Join(services.ErrInternal, err)
	}

	s.pub.PublishForUsers(dto.UUIDs(g.Members), events.ChatDeleted{
		ChatID: chatId,
	})

	return nil
}

func (s *GroupChatService) AddMember(ctx context.Context, chatId, userId uuid.UUID) (*dto.GroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
	}

	err = g.AddMember(domain.UserID(userId))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.GroupMemberAdded{
		ChatID:   chatId,
		MemberID: userId,
	})

	return &gDto, nil
}

func (s *GroupChatService) DeleteMember(ctx context.Context, chatId, memberId uuid.UUID) (*dto.GroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(chatId))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
	}

	err = g.DeleteMember(domain.UserID(memberId))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	gDto := dto.NewGroupChatDTO(g)

	s.pub.PublishForUsers(gDto.Members, events.GroupMemberAdded{
		ChatID:   chatId,
		MemberID: memberId,
	})

	return &gDto, nil
}
