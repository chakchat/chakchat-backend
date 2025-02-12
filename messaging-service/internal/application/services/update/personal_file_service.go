package update

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/google/uuid"
)

type PersonalFileService struct {
	pchatRepo  repository.PersonalChatRepository
	updateRepo repository.UpdateRepository
	files      external.FileStorage
	pub        publish.Publisher
}

func NewPersonalFileService(
	pchatRepo repository.PersonalChatRepository,
	updateRepo repository.UpdateRepository,
	files external.FileStorage,
	pub publish.Publisher,
) *PersonalFileService {
	return &PersonalFileService{
		pchatRepo:  pchatRepo,
		updateRepo: updateRepo,
		files:      files,
		pub:        pub,
	}
}

func (s *PersonalFileService) SendFileMessage(ctx context.Context, req request.SendFileMessage) (*dto.FileMessageDTO, error) {
	chat, err := s.pchatRepo.FindById(ctx, domain.ChatID(req.ChatID))
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

	file, err := s.files.GetById(req.FileID)
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, services.ErrFileNotFound
		}
		return nil, errors.Join(services.ErrInternal, err)
	}
	// For now there is no validation of file.

	domFile := services.NewDomainFileMeta(file)

	msg, err := domain.NewFileMessage(chat, domain.UserID(req.SenderID), &domFile, replyToMessage)
	if err != nil {
		return nil, err
	}

	msg, err = s.updateRepo.CreateFileMessage(ctx, msg)
	if err != nil {
		return nil, errors.Join(services.ErrInternal, err)
	}

	s.pub.PublishForUsers(
		services.GetSecondUserSlice(chat.Members, msg.SenderID),
		events.FileMessageSent{
			ChatID:   uuid.UUID(msg.ChatID),
			UpdateID: int64(msg.UpdateID),
			SenderID: uuid.UUID(msg.SenderID),
			File: events.FileMeta{
				FileId:    msg.File.FileId,
				FileName:  msg.File.FileName,
				MimeType:  msg.File.MimeType,
				FileSize:  msg.File.FileSize,
				FileUrl:   string(msg.File.FileUrl),
				CreatedAt: int64(msg.File.CreatedAt),
			},
			CreatedAt: int64(msg.CreatedAt),
		},
	)

	msgDto := dto.NewFileMessageDTO(msg)
	return &msgDto, nil
}
