package update

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
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

	var payload, iv, keyHash []byte

	if payload, err = base64.StdEncoding.DecodeString(req.Payload); err != nil {
	restapi.SendUnprocessableJSON(c)
		return
	}
	if iv, err = base64.StdEncoding.DecodeString(req.InitializationVector); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}
	if keyHash, err = base64.StdEncoding.DecodeString(req.KeyHash); err != nil {
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

	restapi.SendSuccess(c, generic.FromSecretUpdateDTO(update))
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

	restapi.SendSuccess(c, generic.FromUpdateDeletedDTO(deleted))
}
