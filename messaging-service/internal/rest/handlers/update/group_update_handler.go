package update

import (
	"context"
	"strconv"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupUpdateService interface {
	SendTextMessage(ctx context.Context, req request.SendTextMessage) (*dto.TextMessageDTO, error)
	EditTextMessage(ctx context.Context, req request.EditTextMessage) (*dto.TextMessageDTO, error)
	DeleteMessage(ctx context.Context, req request.DeleteMessage) (*dto.UpdateDeletedDTO, error)
	SendReaction(ctx context.Context, req request.SendReaction) (*dto.ReactionDTO, error)
	DeleteReaction(ctx context.Context, req request.DeleteReaction) (*dto.UpdateDeletedDTO, error)
	ForwardTextMessage(ctx context.Context, req request.ForwardMessage) (*dto.TextMessageDTO, error)
	ForwardFileMessage(ctx context.Context, req request.ForwardMessage) (*dto.FileMessageDTO, error)
}

type GroupUpdateHandler struct {
	service GroupUpdateService
}

func NewGroupUpdateHandler(service GroupUpdateService) *GroupUpdateHandler {
	return &GroupUpdateHandler{
		service: service,
	}
}

func (h *GroupUpdateHandler) SendTextMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
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
		ChatID:         chatID,
		SenderID:       userID,
		Text:           req.Text,
		ReplyToMessage: req.ReplyTo,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromTextMessageDTO(msg))
}

func (h *GroupUpdateHandler) EditTextMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	updateID, err := strconv.ParseInt(c.Param(paramUpdateID), 10, 64)
	if err != nil {
		restapi.SendInvalidUpdateID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	req := struct {
		Text string `json:"text"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
	}

	msg, err := h.service.EditTextMessage(c.Request.Context(), request.EditTextMessage{
		ChatID:    chatID,
		SenderID:  userID,
		MessageID: updateID,
		NewText:   req.Text,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromTextMessageDTO(msg))
}

func (h *GroupUpdateHandler) DeleteMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	updateID, err := strconv.ParseInt(c.Param(paramUpdateID), 10, 64)
	if err != nil {
		restapi.SendInvalidUpdateID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	deleted, err := h.service.DeleteMessage(c.Request.Context(), request.DeleteMessage{
		ChatID:     chatID,
		SenderID:   userID,
		MessageID:  updateID,
		DeleteMode: c.Param(paramDeleteMode),
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromUpdateDeletedDTO(deleted))
}

func (h *GroupUpdateHandler) SendReaction(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	req := struct {
		Reaction  string `json:"reaction"`
		MessageID int64  `json:"message_id"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	reaction, err := h.service.SendReaction(c.Request.Context(), request.SendReaction{
		ChatID:       chatID,
		SenderID:     userID,
		MessageID:    req.MessageID,
		ReactionType: req.Reaction,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromReactionDTO(reaction))
}

func (h *GroupUpdateHandler) DeleteReaction(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	updateID, err := strconv.ParseInt(c.Param(paramUpdateID), 10, 64)
	if err != nil {
		restapi.SendInvalidUpdateID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	deleted, err := h.service.DeleteReaction(c.Request.Context(), request.DeleteReaction{
		ChatID:     chatID,
		SenderID:   userID,
		ReactionID: updateID,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromUpdateDeletedDTO(deleted))
}

func (h *GroupUpdateHandler) ForwardTextMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())
	req := struct {
		Message    int64     `json:"message"`
		FromChatID uuid.UUID `json:"from_chat_id"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	msg, err := h.service.ForwardTextMessage(c.Request.Context(), request.ForwardMessage{
		ToChatID:   chatID,
		SenderID:   userID,
		MessageID:  req.Message,
		FromChatID: req.FromChatID,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromTextMessageDTO(msg))
}

func (h *GroupUpdateHandler) ForwardFileMessage(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())
	req := struct {
		Message    int64     `json:"message"`
		FromChatID uuid.UUID `json:"from_chat_id"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	msg, err := h.service.ForwardFileMessage(c.Request.Context(), request.ForwardMessage{
		ToChatID:   chatID,
		SenderID:   userID,
		MessageID:  req.Message,
		FromChatID: req.FromChatID,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromFileMessageDTO(msg))
}
