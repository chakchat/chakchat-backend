package chat

import (
	"context"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
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
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromSecretPersonalChatDTO(chat))
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
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, generic.FromSecretPersonalChatDTO(chat))
}

func (h *SecretPersonalChatHandler) DeleteChat(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	err = h.service.DeleteChat(c.Request.Context(), request.DeleteChat{
		ChatID:   chatId,
		SenderID: userId,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, struct{}{})

}
