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

type PersonalUpdateService struct {
	txProvider  storage.TxProvider
	pchatRepo   repository.PersonalChatRepository
	updateRepo  repository.UpdateRepository
	chatterRepo repository.ChatterRepository
	pub         publish.Publisher
}

func NewPersonalUpdateService(
	txProvider storage.TxProvider,
	pchatRepo repository.PersonalChatRepository,
	updateRepo repository.UpdateRepository,
	chatterRepo repository.ChatterRepository,
	pub publish.Publisher,
) *PersonalUpdateService {
	return &PersonalUpdateService{
		pchatRepo:   pchatRepo,
		updateRepo:  updateRepo,
		chatterRepo: chatterRepo,
		txProvider:  txProvider,
		pub:         pub,
	}
}

func (s *PersonalUpdateService) SendTextMessage(
	ctx context.Context, req request.SendTextMessage,
) (_ *dto.TextMessageDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.pchatRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	var replyToMessage *domain.Message
	if req.ReplyToMessage != nil {
		replyToMessage, err = s.updateRepo.FindGenericMessage(
			ctx, tx,
			domain.ChatID(req.ChatID),
			domain.UpdateID(*req.ReplyToMessage),
		)
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

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], msg.SenderID, &msg.Update),
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

func (s *PersonalUpdateService) EditTextMessage(
	ctx context.Context, req request.EditTextMessage,
) (_ *dto.TextMessageDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.pchatRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindTextMessage(
		ctx, tx,
		domain.ChatID(req.ChatID),
		domain.UpdateID(req.MessageID),
	)
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

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], msg.Edited.SenderID, &msg.Update),
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

func (s *PersonalUpdateService) DeleteMessage(
	ctx context.Context, req request.DeleteMessage,
) (_ *dto.UpdateDeletedDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.pchatRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindGenericMessage(
		ctx, tx,
		domain.ChatID(req.ChatID),
		domain.UpdateID(req.MessageID),
	)
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

func (s *PersonalUpdateService) SendReaction(
	ctx context.Context, req request.SendReaction,
) (_ *dto.ReactionDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.pchatRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindGenericMessage(
		ctx, tx,
		domain.ChatID(req.ChatID),
		domain.UpdateID(req.MessageID),
	)
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

func (s *PersonalUpdateService) DeleteReaction(
	ctx context.Context, req request.DeleteReaction,
) (_ *dto.UpdateDeletedDTO, err error) {
	tx, err := s.txProvider.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer storage.FinishTx(ctx, tx, &err)

	chat, err := s.pchatRepo.FindById(ctx, tx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	reaction, err := s.updateRepo.FindReaction(
		ctx, tx,
		domain.ChatID(req.ChatID),
		domain.UpdateID(req.ReactionID),
	)
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

	// I do not delete updates because it may cause incosistency.
	// Add triggers and change on delete behavior on foreighn keys before deleting physically

	deleted := reaction.Deleted[len(reaction.Deleted)-1]
	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(chat.Members[:], domain.UserID(req.SenderID), &reaction.Update),
		events.UpdateDeleted{
			ChatID:     uuid.UUID(reaction.ChatID),
			UpdateID:   int64(deleted.UpdateID),
			SenderID:   uuid.UUID(deleted.SenderID),
			DeletedID:  int64(deleted.DeletedID),
			DeleteMode: string(deleted.Mode),
			CreatedAt:  int64(deleted.CreatedAt),
		},
	)

	deletedDto := dto.NewUpdateDeletedDTO(deleted)
	return &deletedDto, nil
}

func (s *PersonalUpdateService) ForwardTextMessage(
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

	toChat, err := s.pchatRepo.FindById(ctx, tx, domain.ChatID(req.ToChatID))
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

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(toChat.Members[:], forwarded.SenderID, &forwarded.Update),
		events.TextMessageSent{
			ChatID:    uuid.UUID(forwarded.ChatID),
			UpdateID:  int64(forwarded.UpdateID),
			SenderID:  uuid.UUID(forwarded.SenderID),
			Text:      forwarded.Text,
			CreatedAt: int64(forwarded.CreatedAt),
		},
	)

	forwardedDto := dto.NewTextMessageDTO(forwarded)
	return &forwardedDto, nil
}

func (s *PersonalUpdateService) ForwardFileMessage(
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

	toChat, err := s.pchatRepo.FindById(ctx, tx, domain.ChatID(req.ToChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	msg, err := s.updateRepo.FindFileMessage(
		ctx, tx, domain.ChatID(req.FromChatID), domain.UpdateID(req.MessageID),
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

	forwarded, err = s.updateRepo.CreateFileMessage(ctx, tx, forwarded)
	if err != nil {
		return nil, err
	}

	s.pub.PublishForUsers(
		services.GetReceivingUpdateMembers(toChat.Members[:], forwarded.SenderID, &forwarded.Update),
		events.FileMessageSent{
			ChatID:   uuid.UUID(forwarded.ChatID),
			UpdateID: int64(forwarded.UpdateID),
			SenderID: uuid.UUID(forwarded.SenderID),
			File: events.FileMeta{
				FileId:    forwarded.File.FileId,
				FileName:  forwarded.File.FileName,
				MimeType:  forwarded.File.MimeType,
				FileSize:  forwarded.File.FileSize,
				FileUrl:   string(forwarded.File.FileURL),
				CreatedAt: int64(forwarded.File.CreatedAt),
			},
			CreatedAt: int64(forwarded.CreatedAt),
		},
	)

	forwardedDto := dto.NewFileMessageDTO(forwarded)
	return &forwardedDto, nil
}
