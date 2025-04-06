package update

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/response"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	queryParamFrom = "from"
	queryParamTo   = "to"
)

type GenericUpdateService interface {
	GetUpdatesRange(context.Context, request.GetUpdatesRange) ([]services.GenericUpdate, error)
	GetUpdate(context.Context, request.GetUpdate) (*services.GenericUpdate, error)
}

type GenericUpdateHandler struct {
	service GenericUpdateService
}

func NewGenericUpdateHandler(service GenericUpdateService) *GenericUpdateHandler {
	return &GenericUpdateHandler{
		service: service,
	}
}

func (h *GenericUpdateHandler) GetUpdatesRange(c *gin.Context) {
	chatID, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	from, err := strconv.ParseInt(c.Query(queryParamFrom), 10, 64)
	if err != nil {
		restapi.SendValidationError(c, []restapi.ErrorDetail{{
			Field:   queryParamFrom,
			Message: "'from' query parameter is required integer",
		}})
	}
	to, err := strconv.ParseInt(c.Query(queryParamTo), 10, 64)
	if err != nil {
		restapi.SendValidationError(c, []restapi.ErrorDetail{{
			Field:   queryParamTo,
			Message: "'to' query parameter is required integer",
		}})
	}
	userID := getUserID(c.Request.Context())

	updates, err := h.service.GetUpdatesRange(c.Request.Context(), request.GetUpdatesRange{
		ChatID:   chatID,
		SenderID: userID,
		From:     from,
		To:       to,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	resp := make([]response.JSONResponse, len(updates))
	for i, up := range updates {
		resp[i] = response.GenericUpdate(&up)
	}
	c.JSON(http.StatusOK, gin.H{
		"updates": resp,
	})
}
