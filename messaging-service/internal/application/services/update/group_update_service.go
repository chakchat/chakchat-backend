package update

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type GroupUpdateService struct {
	groupRepo  repository.GroupChatRepository
	updateRepo repository.UpdateRepository
	txProvider storage.TxProvider
	pub        publish.Publisher
}

func NewGroupUpdateService(
	groupRepo repository.GroupChatRepository,
	updateRepo repository.UpdateRepository,
	txProvider storage.TxProvider,
	pub publish.Publisher,
) *GroupUpdateService {
	return &GroupUpdateService{
		groupRepo:  groupRepo,
		updateRepo: updateRepo,
		txProvider: txProvider,
		pub:        pub,
	}
}

func (s *GroupUpdateService) SendTextMessage(ctx context.Context, req request.SendTextMessage) (*dto.TextMessageDTO, error) {
	chat, err := s.groupRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	var replyToMessage *domain.Message
	if req.ReplyToMessage != nil {
		replyToMessage, err = s.updateRepo.FindGenericMessage(ctx, domain.ChatID(req.ChatID), domain.UpdateID(*req.ReplyToMessage))
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, services.ErrMessageNotFound
			}
			return nil, errors.Join(services.ErrInternal, err)
		}
	}

	msg, err := domain.NewTextMessage(
		chat,
		domain.UserID(req.SenderID),
		req.Text,
		replyToMessage,
	)
	if err != nil {
		return nil, err
	}

	msg, err = s.updateRepo.CreateTextMessage(ctx, msg)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &msg.Update),
		events.TextMessageSent{
			ChatID:    uuid.UUID(msg.ChatID),
			UpdateID:  int64(msg.UpdateID),
			SenderID:  uuid.UUID(msg.SenderID),
			Text:      msg.Text,
			CreatedAt: int64(msg.CreatedAt),
		},
	)

	msgDto := dto.NewTextMessageDTO(msg)
	return &msgDto, nil
}

func (s *GroupUpdateService) EditTextMessage(ctx context.Context, req request.EditTextMessage) (*dto.TextMessageDTO, error) {
	chat, err := s.groupRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	msg, err := s.updateRepo.FindTextMessage(ctx, domain.ChatID(req.ChatID), domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	err = msg.Edit(chat, domain.UserID(req.SenderID), req.NewText)

	if err != nil {
		return nil, err
	}

	err = storage.RunInTx(ctx, s.txProvider, func(ctx context.Context) error {
		msg.Edited, err = s.updateRepo.CreateTextMessageEdited(ctx, msg.Edited)
		if err != nil {
			return err
		}
		msg, err = s.updateRepo.UpdateTextMessage(ctx, msg)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members, msg.Edited.SenderID, &msg.Update),
		events.TextMessageEdited{
			ChatID:    uuid.UUID(msg.ChatID),
			UpdateID:  int64(msg.Edited.UpdateID),
			SenderID:  uuid.UUID(msg.Edited.SenderID),
			MessageID: int64(msg.UpdateID),
			NewText:   msg.Edited.NewText,
			CreatedAt: int64(msg.Edited.CreatedAt),
		},
	)

	msgDto := dto.NewTextMessageDTO(msg)
	return &msgDto, nil
}

func (s *GroupUpdateService) DeleteMessage(ctx context.Context, req request.DeleteMessage) (*dto.UpdateDeletedDTO, error) {
	chat, err := s.groupRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindTextMessage(ctx, domain.ChatID(req.ChatID), domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	err = msg.Delete(chat, domain.UserID(req.SenderID), domain.DeleteMode(req.DeleteMode))

	if err != nil {
		return nil, err
	}

	err = storage.RunInTx(ctx, s.txProvider, func(ctx context.Context) error {
		if msg.DeletedForAll() {
			err := s.updateRepo.DeleteMessage(ctx, msg.ChatID, msg.UpdateID)
			if err != nil {
				return err
			}
		}
		msg.Deleted[len(msg.Deleted)-1], err = s.updateRepo.CreateUpdateDeleted(ctx, msg.Deleted[len(msg.Deleted)-1])
		return err
	})
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	deleted := msg.Deleted[len(msg.Deleted)-1]
	if msg.DeletedForAll() {
		s.pub.PublishForUsers(
			services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &msg.Update),
			events.UpdateDeleted{
				ChatID:     uuid.UUID(msg.ChatID),
				UpdateID:   int64(deleted.UpdateID),
				SenderID:   req.SenderID,
				DeletedID:  req.MessageID,
				DeleteMode: req.DeleteMode,
				CreatedAt:  int64(deleted.CreatedAt),
			},
		)
	}

	deletedDto := dto.NewUpdateDeletedDTO(deleted)
	return &deletedDto, nil
}

func (s *GroupUpdateService) SendReaction(ctx context.Context, req request.SendReaction) (*dto.ReactionDTO, error) {
	chat, err := s.groupRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindGenericMessage(ctx, domain.ChatID(req.ChatID), domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	reaction, err := domain.NewReaction(chat, domain.UserID(req.SenderID), msg, domain.ReactionType(req.ReactionType))
	if err != nil {
		return nil, err
	}

	reaction, err = s.updateRepo.CreateReaction(ctx, reaction)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], reaction.SenderID, &msg.Update),
		events.ReactionSent{
			ChatID:       uuid.UUID(reaction.ChatID),
			UpdateID:     int64(reaction.UpdateID),
			SenderID:     uuid.UUID(reaction.SenderID),
			CreatedAt:    int64(reaction.CreatedAt),
			ReactionType: string(reaction.Type),
		},
	)

	reactionDto := dto.NewReactionDTO(reaction)
	return &reactionDto, nil
}
