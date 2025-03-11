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
)

type FileService struct {
	genChatRepo  repository.GenericChatRepository
	personalSrvc *update.PersonalFileService
	groupSrvc    *update.GroupFileService
}

func NewFileService(
	genChatRepo repository.GenericChatRepository,
	personalSrvc *update.PersonalFileService,
	groupSrvc *update.GroupFileService,
) *FileService {
	return &FileService{
		genChatRepo:  genChatRepo,
		personalSrvc: personalSrvc,
		groupSrvc:    groupSrvc,
	}
}

func (s *FileService) SendFileMessage(ctx context.Context, req request.SendFileMessage) (*dto.FileMessageDTO, error) {
	chatType, err := s.genChatRepo.GetChatType(ctx, domain.ChatID(req.ChatID))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, services.ErrChatNotFound
		}
		return nil, err
	}

	switch chatType {
	case services.ChatTypePersonal:
		return s.personalSrvc.SendFileMessage(ctx, req)
	case services.ChatTypeGroup:
		return s.groupSrvc.SendFileMessage(ctx, req)
	default:
		return nil, errors.Join(services.ErrInvalidChatType,
			fmt.Errorf("cannot send file message in the chat of type: %s", chatType))
	}
}
