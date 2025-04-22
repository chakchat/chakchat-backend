package update

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PersonalFileService interface {
	SendFileMessage(ctx context.Context, req request.SendFileMessage) (*dto.FileMessageDTO, error)
}

type PersonalFileHandler struct {
	service PersonalFileService
}

func NewFileHandler(service PersonalFileService) *PersonalFileHandler {
	return &PersonalFileHandler{
		service: service,
	}
}

func (h *PersonalFileHandler) SendFileMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	req := struct {
		FileID  uuid.UUID `json:"file_id"`
		ReplyTo *int64    `json:"reply_to"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
	}

	msg, err := h.service.SendFileMessage(c.Request.Context(), request.SendFileMessage{
		ChatID:         chatID,
		SenderID:       userID,
		FileID:         req.FileID,
		ReplyToMessage: req.ReplyTo,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromFileMessageDTO(msg))
}
