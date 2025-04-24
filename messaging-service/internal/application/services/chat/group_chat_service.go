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
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain/group"
	"github.com/google/uuid"
)

type GroupChatService struct {
	txProvider storage.TxProvider
	repo       repository.GroupChatRepository
	pub        publish.Publisher
}

func NewGroupChatService(
	txProvider storage.TxProvider, repo repository.GroupChatRepository, pub publish.Publisher,
) *GroupChatService {
	return &GroupChatService{
		repo:       repo,
		pub:        pub,
		txProvider: txProvider,
	}
}

func (s *GroupChatService) CreateGroup(ctx context.Context, req request.CreateGroup) (_ *dto.GroupChatDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)
	members := make([]domain.UserID, len(req.Members))
	for i, m := range req.Members {
		members[i] = domain.UserID(m)
	}

	g, err := group.NewGroupChat(domain.UserID(req.SenderID), members, req.Name)

	if err != nil {
		return nil, err
	}

	g, err = s.repo.Create(ctx, tx, g)
	if err != nil {
		return nil, err
	}

	gDto := dto.NewGroupChatDTO(g)

	err = s.pub.PublishForReceivers(
		ctx,
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeChatCreated,
		events.ChatCreated{
			SenderID: req.SenderID,
			Chat:     generic.FromGroupChatDTO(&gDto),
		},
	)
	if err != nil {
		return nil, err
	}

	return &gDto, nil
}

func (s *GroupChatService) UpdateGroupInfo(ctx context.Context, req request.UpdateGroupInfo) (_ *dto.GroupChatDTO, err error) {
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

	gDto := dto.NewGroupChatDTO(g)

	err = s.pub.PublishForReceivers(
		ctx,
		services.GetReceivingMembers(g.Members, g.Admin),
		events.TypeGroupInfoUpdated,
		events.GroupInfoUpdated{
			SenderID:    req.SenderID,
			ChatID:      gDto.ID,
			Name:        gDto.Name,
			Description: gDto.Description,
			GroupPhoto:  string(g.GroupPhoto),
		},
	)
	if err != nil {
		return nil, err
	}

	return &gDto, nil
}

func (s *GroupChatService) DeleteGroup(ctx context.Context, req request.DeleteChat) (err error) {
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

	err = s.pub.PublishForReceivers(
		ctx,
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeChatDeleted,
		events.ChatDeleted{
			SenderID: req.SenderID,
			ChatID:   req.ChatID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *GroupChatService) AddMember(ctx context.Context, req request.AddMember) (_ *dto.GroupChatDTO, err error) {
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

	gDto := dto.NewGroupChatDTO(g)

	err = s.pub.PublishForReceivers(
		ctx,
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeGroupMembersAdded,
		events.GroupMemberAdded{
			SenderID: req.SenderID,
			ChatID:   req.ChatID,
			Members:  []uuid.UUID{req.MemberID},
		},
	)
	if err != nil {
		return nil, err
	}

	return &gDto, nil
}

func (s *GroupChatService) DeleteMember(ctx context.Context, req request.DeleteMember) (_ *dto.GroupChatDTO, err error) {
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

	gDto := dto.NewGroupChatDTO(g)

	err = s.pub.PublishForReceivers(
		ctx,
		services.GetReceivingMembers(g.Members, domain.UserID(req.SenderID)),
		events.TypeGroupMembersRemoved,
		events.GroupMembersRemoved{
			SenderID: req.SenderID,
			ChatID:   req.ChatID,
			Members:  []uuid.UUID{req.MemberID},
		},
	)
	if err != nil {
		return nil, err
	}

	return &gDto, nil
}
