package chat

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/secgroup"
	"github.com/google/uuid"
)

type SecretGroupChatService struct {
	txProvider storage.TxProvider
	repo       repository.SecretGroupChatRepository
	pub        publish.Publisher
}

func NewSecretGroupChatService(
	txProvider storage.TxProvider, repo repository.SecretGroupChatRepository, pub publish.Publisher,
) *SecretGroupChatService {
	return &SecretGroupChatService{
		txProvider: txProvider,
		repo:       repo,
		pub:        pub,
	}
}

func (s *SecretGroupChatService) CreateGroup(
	ctx context.Context, req request.CreateSecretGroup,
) (_ *dto.SecretGroupChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	g, err := secgroup.NewSecretGroupChat(domain.UserID(req.SenderID), members, req.Name)

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Create(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeChatCreated,
		events.ChatCreated{
			SenderID: req.SenderID,
			Chat:     generic.FromSecretGroupChatDTO(&gDto),
		},
	)

	return &gDto, nil
}

func (s *SecretGroupChatService) UpdateGroupInfo(
	ctx context.Context, req request.UpdateSecretGroupInfo,
) (_ *dto.SecretGroupChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	g, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
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

	g, err = s.repo.Update(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeGroupInfoUpdated,
		events.GroupInfoUpdated{
			SenderID:    req.SenderID,
			ChatID:      gDto.ID,
			Name:        gDto.Name,
			Description: gDto.Description,
			GroupPhoto:  string(g.GroupPhoto),
		},
	)

	return &gDto, nil
}

func (s *SecretGroupChatService) DeleteGroup(ctx context.Context, req request.DeleteChat) (err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return err
	}
	defer storage.FinishTx(ctx, tx, &err)

	g, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
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

	if err := s.repo.Delete(ctx, tx, g.ID); err != nil {
		return err
	}

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeChatDeleted,
		events.ChatDeleted{
			SenderID: req.SenderID,
			ChatID:   req.ChatID,
		},
	)

	return nil
}
func (s *SecretGroupChatService) AddMember(
	ctx context.Context, req request.AddMember,
) (_ *dto.SecretGroupChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	g, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
	}

	err = g.AddMember(domain.UserID(req.SenderID), domain.UserID(req.MemberID))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeGroupMembersAdded,
		events.GroupMemberAdded{
			SenderID: req.SenderID,
			ChatID:   req.ChatID,
			Members:  []uuid.UUID{req.MemberID},
		},
	)

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}

func (s *SecretGroupChatService) DeleteMember(
	ctx context.Context, req request.DeleteMember,
) (_ *dto.SecretGroupChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	g, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
	}

	err = g.DeleteMember(domain.UserID(req.SenderID), domain.UserID(req.MemberID))

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Update(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewSecretGroupChatDTO(g)

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeGroupMembersRemoved,
		events.GroupMemberAdded{
			SenderID: req.SenderID,
			ChatID:   req.ChatID,
			Members:  []uuid.UUID{req.MemberID},
		},
	)

	return &gDto, nil
}

func (s *SecretGroupChatService) SetExpiration(
	ctx context.Context, req request.SetExpiration,
) (_ *dto.SecretGroupChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	g, err := s.repo.FindById(ctx, tx, domain.ChatID(req.ChatID))
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

	g, err = s.repo.Update(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForReceivers(
		services.GetReceivingMembers(g.Members[:], domain.UserID(req.SenderID)),
		events.TypeChatExpirationSet,
		events.ExpirationSet{
			SenderID:   req.SenderID,
			ChatID:     req.ChatID,
			Expiration: req.Expiration,
		},
	)

	gDto := dto.NewSecretGroupChatDTO(g)
	return &gDto, nil
}
