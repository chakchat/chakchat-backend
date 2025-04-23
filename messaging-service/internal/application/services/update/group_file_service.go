package update

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/external"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/publish/events"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GroupFileService struct {
	txProvider storage.TxProvider
	groupRepo  repository.GroupChatRepository
	updateRepo repository.UpdateRepository
	files      external.FileStorage
	pub        publish.Publisher
}

func NewGroupFileService(
	txProvider storage.TxProvider,
	groupRepo repository.GroupChatRepository,
	updateRepo repository.UpdateRepository,
	files external.FileStorage,
	pub publish.Publisher,
) *GroupFileService {
	return &GroupFileService{
		groupRepo:  groupRepo,
		updateRepo: updateRepo,
		files:      files,
		pub:        pub,
		txProvider: txProvider,
	}
}

func (s *GroupFileService) SendFileMessage(
	ctx context.Context, req request.SendFileMessage,
) (_ *dto.FileMessageDTO, err error) {
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

	file, err := s.files.GetById(ctx, req.FileID)
	if err != nil {
		if errors.Is(err, external.ErrFileNotFound) {
			return nil, services.ErrFileNotFound
		}
		return nil, err
	}
	// For now there is no validation of file.

	domFile := services.NewDomainFileMeta(file)

	msg, err := domain.NewFileMessage(chat, domain.UserID(req.SenderID), &domFile, replyToMessage)
	if err != nil {
		return nil, err
	}

	msg, err = s.updateRepo.CreateFileMessage(ctx, tx, msg)
	if err != nil {
		return nil, err
	}

	msgDto := dto.NewFileMessageDTO(msg)

	s.pub.PublishForReceivers(
		services.GetReceivingUpdateMembers(chat.Members[:], msg.SenderID, &msg.Update),
		events.TypeUpdate,
		generic.FromFileMessageDTO(&msgDto),
	)

	return &msgDto, nil
}
