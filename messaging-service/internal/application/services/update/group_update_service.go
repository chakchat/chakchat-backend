package update

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
)

type GroupUpdateService struct {
	txProvider  storage.TxProvider
	groupRepo   repository.GroupChatRepository
	updateRepo  repository.UpdateRepository
	chatterRepo repository.ChatterRepository
	pub         publish.Publisher
}

func NewGroupUpdateService(
	txProvider storage.TxProvider,
	groupRepo repository.GroupChatRepository,
	updateRepo repository.UpdateRepository,
	chatterRepo repository.ChatterRepository,
	pub publish.Publisher,
) *GroupUpdateService {
	return &GroupUpdateService{
		groupRepo:   groupRepo,
		updateRepo:  updateRepo,
		chatterRepo: chatterRepo,
		txProvider:  txProvider,
		pub:         pub,
	}
}

func (s *GroupUpdateService) SendTextMessage(
	ctx context.Context, req request.SendTextMessage,
) (_ *dto.TextMessageDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.groupRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	var replyToMessage *domain.Message
	if req.ReplyToMessage != nil {
		replyToMessage, err = s.updateRepo.FindGenericMessage(ctx, tx, domain.ChatID(req.ChatID), domain.UpdateID(*req.ReplyToMessage))
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, services.ErrMessageNotFound
			}
			return nil, err
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

	msg, err = s.updateRepo.CreateTextMessage(ctx, tx, msg)
	if err != nil {
		return nil, err
	}

	msgDto := dto.NewTextMessageDTO(msg)

	s.pub.PublishForReceivers(
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &msg.Update),
		events.TypeUpdate,
		generic.FromTextMessageDTO(&msgDto),
	)

	return &msgDto, nil
}

func (s *GroupUpdateService) EditTextMessage(
	ctx context.Context, req request.EditTextMessage,
) (_ *dto.TextMessageDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.groupRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindTextMessage(ctx, tx, domain.ChatID(req.ChatID), domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, err
	}

	err = msg.Edit(chat, domain.UserID(req.SenderID), req.NewText)

	if err != nil {
		return nil, err
	}

	msg.Edited, err = s.updateRepo.CreateTextMessageEdited(ctx, tx, msg.Edited)
	if err != nil {
		return nil, err
	}
	msg, err = s.updateRepo.UpdateTextMessage(ctx, tx, msg)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	msgDto := dto.NewTextMessageDTO(msg)

	s.pub.PublishForReceivers(
		services.GetReceivingUpdateMembers(chat.Members, msg.Edited.SenderID, &msg.Update),
		events.TypeUpdate,
		generic.FromTextMessageEditedDTO(msgDto.Edited),
	)

	return &msgDto, nil
}

func (s *GroupUpdateService) DeleteMessage(
	ctx context.Context, req request.DeleteMessage,
) (_ *dto.UpdateDeletedDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.groupRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindTextMessage(ctx, tx, domain.ChatID(req.ChatID), domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, err
	}

	deleteMode, err := domain.NewDeleteMode(req.DeleteMode)
	if err != nil {
		return nil, err
	}

	err = msg.Delete(chat, domain.UserID(req.SenderID), deleteMode)

	if err != nil {
		return nil, err
	}
	msg.Deleted[len(msg.Deleted)-1], err = s.updateRepo.CreateUpdateDeleted(
		ctx, tx, msg.Deleted[len(msg.Deleted)-1],
	)
	if err != nil {
		return nil, err
	}

	// I do not delete updates because it may cause incosistency.
	// Add triggers and change on delete behavior on foreighn keys before deleting physically

	deleted := msg.Deleted[len(msg.Deleted)-1]
	deletedDto := dto.NewUpdateDeletedDTO(deleted)

	if msg.DeletedForAll() {
		s.pub.PublishForReceivers(
			services.GetReceivingUpdateMembers(chat.Members, domain.UserID(req.SenderID), &msg.Update),
			events.TypeUpdate,
			generic.FromUpdateDeletedDTO(&deletedDto),
		)
	}

	return &deletedDto, nil
}

func (s *GroupUpdateService) SendReaction(
	ctx context.Context, req request.SendReaction,
) (_ *dto.ReactionDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.groupRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindGenericMessage(ctx, tx, domain.ChatID(req.ChatID), domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, err
	}

	reaction, err := domain.NewReaction(chat, domain.UserID(req.SenderID), msg, domain.ReactionType(req.ReactionType))
	if err != nil {
		return nil, err
	}

	reaction, err = s.updateRepo.CreateReaction(ctx, tx, reaction)
	if err != nil {
		return nil, err
	}

	reactionDto := dto.NewReactionDTO(reaction)

	s.pub.PublishForReceivers(
		services.GetReceivingUpdateMembers(chat.Members, reaction.SenderID, &msg.Update),
		events.TypeUpdate,
		generic.FromReactionDTO(&reactionDto),
	)

	return &reactionDto, nil
}

func (s *GroupUpdateService) DeleteReaction(
	ctx context.Context, req request.DeleteReaction,
) (_ *dto.UpdateDeletedDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.groupRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	reaction, err := s.updateRepo.FindReaction(ctx, tx, domain.ChatID(req.ChatID), domain.UpdateID(req.ReactionID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrReactionNotFound
		}
		return nil, err
	}

	err = reaction.Delete(chat, domain.UserID(req.SenderID))
	if err != nil {
		return nil, err
	}

	reaction.Deleted[len(reaction.Deleted)-1], err = s.updateRepo.CreateUpdateDeleted(
		ctx, tx, reaction.Deleted[len(reaction.Deleted)-1],
	)
	if err != nil {
		return nil, err
	}

	// For now reaction is always deleted for all users. And no `if reaction.DeletedForAll() {...}` check is performed.
	err = s.updateRepo.DeleteUpdate(ctx, tx, reaction.ChatID, reaction.UpdateID)
	if err != nil {
		return nil, err
	}

	deleted := reaction.Deleted[len(reaction.Deleted)-1]
	deletedDto := dto.NewUpdateDeletedDTO(deleted)

	s.pub.PublishForReceivers(
		services.GetReceivingUpdateMembers(chat.Members, domain.UserID(req.SenderID), &reaction.Update),
		events.TypeChatDeleted,
		generic.FromUpdateDeletedDTO(&deletedDto),
	)

	return &deletedDto, nil
}

func (s *GroupUpdateService) ForwardTextMessage(
	ctx context.Context, req request.ForwardMessage,
) (_ *dto.TextMessageDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	fromChat, err := s.chatterRepo.FindChatter(ctx, tx, domain.ChatID(req.FromChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	toChat, err := s.groupRepo.FindById(ctx, tx, domain.ChatID(req.ToChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindTextMessage(
		ctx, tx,
		domain.ChatID(req.FromChatID),
		domain.UpdateID(req.MessageID),
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, err
	}

	forwarded, err := msg.Forward(fromChat, domain.UserID(req.SenderID), toChat)
	if err != nil {
		return nil, err
	}

	forwarded, err = s.updateRepo.CreateTextMessage(ctx, tx, forwarded)
	if err != nil {
		return nil, err
	}

	forwardedDto := dto.NewTextMessageDTO(forwarded)

	s.pub.PublishForReceivers(
		services.GetReceivingUpdateMembers(toChat.Members[:], forwarded.SenderID, &forwarded.Update),
		events.TypeUpdate,
		generic.FromTextMessageDTO(&forwardedDto),
	)

	return &forwardedDto, nil
}

func (s *GroupUpdateService) ForwardFileMessage(
	ctx context.Context, req request.ForwardMessage,
) (_ *dto.FileMessageDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	fromChat, err := s.chatterRepo.FindChatter(ctx, tx, domain.ChatID(req.FromChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	toChat, err := s.groupRepo.FindById(ctx, tx, domain.ChatID(req.ToChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindFileMessage(ctx, tx, domain.ChatID(req.FromChatID), domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, err
	}

	forwarded, err := msg.Forward(fromChat, domain.UserID(req.SenderID), toChat)
	if err != nil {
		return nil, err
	}

	forwarded, err = s.updateRepo.CreateFileMessage(ctx, tx, forwarded)
	if err != nil {
		return nil, err
	}

	forwardedDto := dto.NewFileMessageDTO(forwarded)

	s.pub.PublishForReceivers(
		services.GetReceivingUpdateMembers(toChat.Members[:], forwarded.SenderID, &forwarded.Update),
		events.TypeUpdate,
		generic.FromFileMessageDTO(&forwardedDto),
	)

	return &forwardedDto, nil
}
