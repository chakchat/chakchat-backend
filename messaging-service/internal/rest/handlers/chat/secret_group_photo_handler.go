package chat

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecretGroupPhotoService interface {
	UpdatePhoto(ctx context.Context, req request.UpdateGroupPhoto) (*dto.GroupChatDTO, error)
	DeletePhoto(ctx context.Context, req request.DeleteGroupPhoto) (*dto.GroupChatDTO, error)
}

type SecretGroupPhotoHandler struct {
	service SecretGroupPhotoService
}

func NewSecretGroupPhotoHandler(service SecretGroupPhotoService) *SecretGroupPhotoHandler {
	return &SecretGroupPhotoHandler{
		service: service,
	}
}

func (h *SecretGroupPhotoHandler) UpdatePhoto(c *gin.Context) {
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
		resp := errmap.Map(err)
		c.JSON(resp.Code, resp.Body)
		return
	}

	restapi.SendSuccess(c, newGroupResponse(group))
}

func (h *SecretGroupPhotoHandler) DeletePhoto(c *gin.Context) {
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
		resp := errmap.Map(err)
		c.JSON(resp.Code, resp.Body)
		return
	}

	restapi.SendSuccess(c, newGroupResponse(group))
}
