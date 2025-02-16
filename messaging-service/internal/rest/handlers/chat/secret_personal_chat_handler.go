package chat

import (
	"context"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecretPersonalChatService interface {
	CreateChat(ctx context.Context, req request.CreateSecretPersonalChat) (*dto.SecretPersonalChatDTO, error)
	SetExpiration(ctx context.Context, req request.SetExpiration) (*dto.SecretPersonalChatDTO, error)
	DeleteChat(ctx context.Context, req request.DeleteChat) error
}

type SecretPersonalChatHandler struct {
	service SecretPersonalChatService
}

func NewSecretPersonalChatHandler(service SecretPersonalChatService) *SecretPersonalChatHandler {
	return &SecretPersonalChatHandler{
		service: service,
	}
}

func (h *SecretPersonalChatHandler) CreateChat(c *gin.Context) {
	userId := getUserID(c.Request.Context())

	req := struct {
		MemberID uuid.UUID `json:"member_id"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	chat, err := h.service.CreateChat(c.Request.Context(), request.CreateSecretPersonalChat{
		SenderID: userId,
		MemberID: req.MemberID,
	})
	if err != nil {
		resp := errmap.Map(err)
		c.JSON(resp.Code, resp.Body)
		return
	}

	restapi.SendSuccess(c, newSecretPersonalChatResponse(chat))
}

func (h *SecretPersonalChatHandler) SetExpiration(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	req := struct {
		Expiration *time.Duration `json:"expiration"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.Error(err)
		return
	}

	chat, err := h.service.SetExpiration(c.Request.Context(), request.SetExpiration{
		ChatID:     chatId,
		SenderID:   userId,
		Expiration: req.Expiration,
	})
	if err != nil {
		resp := errmap.Map(err)
		c.JSON(resp.Code, resp.Body)
		return
	}

	restapi.SendSuccess(c, newSecretPersonalChatResponse(chat))
}

type secretPersonalChatResponse struct {
	ID      uuid.UUID    `json:"chat_id"`
	Members [2]uuid.UUID `json:"members"`

	Expiration *time.Duration `json:"expiration"`
	CreatedAt  int64          `json:"created_at"`
}

func newSecretPersonalChatResponse(dto *dto.SecretPersonalChatDTO) secretPersonalChatResponse {
	return secretPersonalChatResponse{
		ID:         dto.ID,
		Members:    dto.Members,
		Expiration: dto.Expiration,
		CreatedAt:  dto.CreatedAt,
	}
}
