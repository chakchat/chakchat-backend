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
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secgroup"
	"github.com/google/uuid"
)

type SecretGroupChatService struct {
	repo repository.SecretGroupChatRepository
	pub  publish.Publisher
}

func NewSecretGroupChatService(repo repository.SecretGroupChatRepository, pub publish.Publisher) *SecretGroupChatService {
	return &SecretGroupChatService{
		repo: repo,
		pub:  pub,
	}
}

func (s *SecretGroupChatService) CreateGroup(ctx context.Context, req request.CreateSecretGroup) (*dto.SecretGroupChatDTO, error) {
	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	g, err := secgroup.NewSecretGroupChat(domain.UserID(req.SenderID), members, req.Name)

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Create(ctx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.ChatCreated{
			ChatID:   uuid.UUID(gDto.ID),
			ChatType: events.ChatTypeSecretGroup,
		},
	)

	return &gDto, nil
}

func (s *SecretGroupChatService) UpdateGroupInfo(ctx context.Context, req request.UpdateSecretGroupInfo) (*dto.SecretGroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	err = g.UpdateInfo(domain.UserID(req.SenderID), req.Name, req.Description)

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.GroupInfoUpdated{
			ChatID:        gDto.ID,
			Name:          gDto.Name,
			Description:   gDto.Description,
			GroupPhotoURL: string(g.GroupPhoto),
		},
	)

	return &gDto, nil
}

func (s *SecretGroupChatService) DeleteGroup(ctx context.Context, req request.DeleteChat) error {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return services.ErrChatNotFound
		}
		return err
	}

	err = g.Delete(domain.UserID(req.SenderID))
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, g.ID); err != nil {
		return err
	}

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.ChatDeleted{
			ChatID: req.ChatID,
		},
	)

	return nil
}
func (s *SecretGroupChatService) AddMember(ctx context.Context, req request.AddMember) (*dto.SecretGroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
	}

	err = g.AddMember(domain.UserID(req.SenderID), domain.UserID(req.MemberID))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.GroupMemberAdded{
			ChatID:   req.ChatID,
			MemberID: req.MemberID,
		},
	)

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}

func (s *SecretGroupChatService) DeleteMember(ctx context.Context, req request.DeleteMember) (*dto.SecretGroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
	}

	err = g.DeleteMember(domain.UserID(req.SenderID), domain.UserID(req.MemberID))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.GroupMemberAdded{
			ChatID:   req.ChatID,
			MemberID: req.MemberID,
		},
	)

	return &gDto, nil
}

func (s *SecretGroupChatService) SetExpiration(ctx context.Context, req request.SetExpiration) (*dto.SecretGroupChatDTO, error) {
	g, err := s.repo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	err = g.SetExpiration(domain.UserID(req.SenderID), req.Expiration)
	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, g)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForUsers(
		services.GetReceivingMembers(g.Members[:], domain.UserID(req.SenderID)),
		events.ExpirationSet{
			ChatID:     req.ChatID,
			SenderID:   req.SenderID,
			Expiration: req.Expiration,
		},
	)

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}
