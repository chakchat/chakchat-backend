package updatex

import (
	"context"
	"errors"
	"fmt"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/update"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type UpdateService struct {
	genChatRepo  repository.GenericChatRepository
	personalSrvc *update.PersonalUpdateService
	groupSrvc    *update.GroupUpdateService
}

func NewUpdateService(
	genChatRepo repository.GenericChatRepository,
	personalSrvc *update.PersonalUpdateService,
	groupSrvc *update.GroupUpdateService,
) *UpdateService {
	return &UpdateService{
		genChatRepo:  genChatRepo,
		personalSrvc: personalSrvc,
		groupSrvc:    groupSrvc,
	}
}

func (s *UpdateService) SendTextMessage(ctx context.Context, req request.SendTextMessage) (*dto.TextMessageDTO, error) {
	chatType, err := s.getChatType(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.SendTextMessage(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.SendTextMessage(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot send text message in the chat of type: %s", chatType))
	}
}

func (s *UpdateService) EditTextMessage(ctx context.Context, req request.EditTextMessage) (*dto.TextMessageDTO, error) {
	chatType, err := s.getChatType(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.EditTextMessage(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.EditTextMessage(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot edit text message in the chat of type: %s", chatType))
	}
}

func (s *UpdateService) DeleteMessage(ctx context.Context, req request.DeleteMessage) (*dto.UpdateDeletedDTO, error) {
	chatType, err := s.getChatType(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.DeleteMessage(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.DeleteMessage(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot delete message in the chat of type: %s", chatType))
	}
}

func (s *UpdateService) SendReaction(ctx context.Context, req request.SendReaction) (*dto.ReactionDTO, error) {
	chatType, err := s.getChatType(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.SendReaction(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.SendReaction(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot send reaction in the chat of type: %s", chatType))
	}
}

func (s *UpdateService) DeleteReaction(ctx context.Context, req request.DeleteReaction) (*dto.UpdateDeletedDTO, error) {
	chatType, err := s.getChatType(ctx, req.ChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.DeleteReaction(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.DeleteReaction(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot delete reaction in the chat of type: %s", chatType))
	}
}

func (s *UpdateService) ForwardTextMessage(ctx context.Context, req request.ForwardMessage) (*dto.TextMessageDTO, error) {
	chatType, err := s.getChatType(ctx, req.ToChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.ForwardTextMessage(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.ForwardTextMessage(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot forward text message in the chat of type: %s", chatType))
	}
}

func (s *UpdateService) ForwardFileMessage(ctx context.Context, req request.ForwardMessage) (*dto.FileMessageDTO, error) {
	chatType, err := s.getChatType(ctx, req.ToChatID)
	if err != nil {
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.ForwardFileMessage(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.ForwardFileMessage(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot forward file message in the chat of type: %s", chatType))
	}
}

func (s *UpdateService) getChatType(ctx context.Context, id uuid.UUID) (string, error) {
	chatType, err := s.genChatRepo.GetChatType(ctx, domain.ChatID(id))
	if errors.Is(err, repository.ErrNotFound) {
		return "", services.ErrChatNotFound
	}
	return chatType, err
}
