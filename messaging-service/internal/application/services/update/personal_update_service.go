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
	pchatRepo  repository.PersonalChatRepository
	updateRepo repository.UpdateRepository
	txProvider storage.TxProvider
	pub        publish.Publisher
}

func NewPersonalUpdateService(
	pchatRepo repository.PersonalChatRepository,
	updateRepo repository.UpdateRepository,
	transactioner storage.TxProvider,
	pub publish.Publisher,
) *PersonalUpdateService {
	return &PersonalUpdateService{
		pchatRepo:  pchatRepo,
		updateRepo: updateRepo,
		txProvider: transactioner,
		pub:        pub,
	}
}

func (s *PersonalUpdateService) SendTextMessage(ctx context.Context, req request.SendTextMessage) (*dto.TextMessageDTO, error) {
	chat, err := s.pchatRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	var replyToMessage *domain.Message
	if req.ReplyToMessage != nil {
		replyToMessage, err = s.updateRepo.FindGenericMessage(ctx, domain.UpdateID(*req.ReplyToMessage))
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

	switch {
	case errors.Is(err, domain.ErrUserNotMember):
		return nil, services.ErrUserNotMember
	case errors.Is(err, domain.ErrChatBlocked):
		return nil, services.ErrChatBlocked
	case errors.Is(err, domain.ErrUpdateDeleted):
		return nil, services.ErrUpdateDeleted
	case errors.Is(err, domain.ErrUpdateNotFromChat):
		return nil, services.ErrUpdateNotFromChat
	case errors.Is(err, domain.ErrTextEmpty):
		return nil, services.ErrTextEmpty
	case errors.Is(err, domain.ErrTooMuchTextRunes):
		return nil, services.ErrTooMuchTextRunes
	case err != nil:
		return nil, errors.Join(services.ErrInternal, err)
	}

	msg, err = s.updateRepo.CreateTextMessage(ctx, msg)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	s.pub.PublishForUsers(
		services.GetSecondUserSlice(chat.Members, msg.SenderID),
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

func (s *PersonalUpdateService) EditTextMessage(ctx context.Context, req request.EditTextMessage) (*dto.TextMessageDTO, error) {
	chat, err := s.pchatRepo.FindById(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	msg, err := s.updateRepo.FindTextMessage(ctx, domain.UpdateID(req.MessageID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrMessageNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}

	err = msg.Edit(chat, domain.UserID(req.SenderID), req.NewText)

	switch {
	case errors.Is(err, domain.ErrUserNotMember):
		return nil, services.ErrUserNotMember
	case errors.Is(err, domain.ErrChatBlocked):
		return nil, services.ErrChatBlocked
	case errors.Is(err, domain.ErrUpdateNotFromChat):
		return nil, services.ErrUpdateNotFromChat
	case errors.Is(err, domain.ErrUserNotSender):
		return nil, services.ErrUserNotSender
	case errors.Is(err, domain.ErrUpdateDeleted):
		return nil, domain.ErrUpdateDeleted
	case errors.Is(err, domain.ErrTextEmpty):
		return nil, services.ErrTextEmpty
	case errors.Is(err, domain.ErrTooMuchTextRunes):
		return nil, services.ErrTooMuchTextRunes
	case err != nil:
		return nil, errors.Join(services.ErrInternal, err)
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
		services.GetSecondUserSlice(chat.Members, msg.Edited.SenderID),
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
