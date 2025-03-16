package chat

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

type GroupPhotoService interface {
	UpdatePhoto(ctx context.Context, req request.UpdateGroupPhoto) (*dto.GroupChatDTO, error)
	DeletePhoto(ctx context.Context, req request.DeleteGroupPhoto) (*dto.GroupChatDTO, error)
}

type GroupPhotoHandler struct {
	service GroupPhotoService
}

func NewGroupPhotoHandler(service GroupPhotoService) *GroupPhotoHandler {
	return &GroupPhotoHandler{
		service: service,
	}
}

func (h *GroupPhotoHandler) UpdatePhoto(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	req := struct {
		PhotoID uuid.UUID `json:"photo_id"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	group, err := h.service.UpdatePhoto(c.Request.Context(), request.UpdateGroupPhoto{
		ChatID:   chatId,
		SenderID: userId,
		FileID:   req.PhotoID,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.GroupChat(group))
}

func (h *GroupPhotoHandler) DeletePhoto(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	req := struct {
		PhotoID uuid.UUID `json:"photo_id"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	group, err := h.service.DeletePhoto(c.Request.Context(), request.DeleteGroupPhoto{
		ChatID:   chatId,
		SenderID: userId,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.GroupChat(group))
}
