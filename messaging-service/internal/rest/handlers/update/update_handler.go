package update

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/response"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateService interface {
	SendTextMessage(ctx context.Context, req request.SendTextMessage) (*dto.TextMessageDTO, error)
	EditTextMessage(ctx context.Context, req request.EditTextMessage) (*dto.TextMessageDTO, error)
	DeleteMessage(ctx context.Context, req request.DeleteMessage) (*dto.UpdateDeletedDTO, error)
	SendReaction(ctx context.Context, req request.SendReaction) (*dto.ReactionDTO, error)
	DeleteReaction(ctx context.Context, req request.DeleteReaction) (*dto.UpdateDeletedDTO, error)
	ForwardTextMessage(ctx context.Context, req request.ForwardMessage) (*dto.TextMessageDTO, error)
	ForwardFileMessage(ctx context.Context, req request.ForwardMessage) (*dto.FileMessageDTO, error)
}

const paramChatID = "chatId"

type UpdateHandler struct {
	service UpdateService
}

func NewUpdateHandler(service UpdateService) *UpdateHandler {
	return &UpdateHandler{
		service: service,
	}
}

func (h *UpdateHandler) SendTextMessage(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	req := struct {
		Text    string `json:"text"`
		ReplyTo *int64 `json:"reply_to"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
	}

	msg, err := h.service.SendTextMessage(c.Request.Context(), request.SendTextMessage{
		ChatID:         chatId,
		SenderID:       userID,
		Text:           req.Text,
		ReplyToMessage: req.ReplyTo,
	})
	if err != nil {
		errmap.Respond(c, err)
	}

	restapi.SendSuccess(c, response.TextMessage(msg))
}
