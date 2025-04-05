package update

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/response"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecretPersonalUpdateService interface {
	SendSecretUpdate(context.Context, request.SendSecretUpdate) (*dto.SecretUpdateDTO, error)
	DeleteSecretUpdate(context.Context, request.DeleteSecretUpdate) (*dto.UpdateDeletedDTO, error)
}

type SecretPersonalUpdateHandler struct {
	service SecretPersonalUpdateService
}

func NewSecretPersonalUpdateHandler(service SecretPersonalUpdateService) *SecretPersonalUpdateHandler {
	return &SecretPersonalUpdateHandler{
		service: service,
	}
}

func (h *SecretPersonalUpdateHandler) SendSecretUpdate(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	req := struct {
		Payload              string `json:"payload"`
		InitializationVector string `json:"initialization_vector"`
		KeyHash              string `json:"key_hash"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	var (
		payload = make([]byte, 0, len(req.Payload))
		iv      = make([]byte, 0, len(req.InitializationVector))
		keyHash = make([]byte, 0, len(req.KeyHash))
	)
	if _, err := base64.StdEncoding.Decode(payload, []byte(req.Payload)); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}
	if _, err := base64.StdEncoding.Decode(iv, []byte(req.InitializationVector)); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}
	if _, err := base64.StdEncoding.Decode(keyHash, []byte(req.KeyHash)); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	update, err := h.service.SendSecretUpdate(c.Request.Context(), request.SendSecretUpdate{
		ChatID:               chatID,
		SenderID:             userID,
		Payload:              payload,
		InitializationVector: iv,
		KeyHash:              keyHash,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.SecretUpdate(update))
}

func (h *SecretPersonalUpdateHandler) DeleteSecretUpdate(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())
	updateID, err := strconv.ParseInt(c.Param(paramUpdateID), 10, 64)
	if err != nil {
		restapi.SendInvalidUpdateID(c)
		return
	}

	deleted, err := h.service.DeleteSecretUpdate(c.Request.Context(), request.DeleteSecretUpdate{
		ChatID:         chatID,
		SenderID:       userID,
		SecretUpdateID: updateID,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.UpdateDeleted(deleted))
}
