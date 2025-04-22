package chat

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	queryParamLastUpdateID = "enableLastUpdateID"
	queryParamPreviewCount = "preview"
)

type GenericChatService interface {
	GetByMemberID(
		ctx context.Context, memberID uuid.UUID, opts ...request.GetChatOption,
	) ([]generic.Chat, error)
	GetByChatID(
		ctx context.Context, senderID, chatID uuid.UUID, opts ...request.GetChatOption,
	) (*generic.Chat, error)
}

type GenericChatHandler struct {
	service GenericChatService
}

func NewGenericChatHandler(service GenericChatService) *GenericChatHandler {
	return &GenericChatHandler{
		service: service,
	}
}

func (h *GenericChatHandler) GetAllChats(c *gin.Context) {
	userID := getUserID(c.Request.Context())

	var opts []request.GetChatOption
	if enableLastUpdateID := c.Query(queryParamLastUpdateID); enableLastUpdateID != "" {
		switch enableLastUpdateID {
		case "true":
			opts = append(opts, request.WithChatLastUpdateID())
		case "false":
		default:
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    "invalid_query_param",
				ErrorMessage: "Invalid query parameter",
				ErrorDetails: []restapi.ErrorDetail{{
					Field:   queryParamLastUpdateID,
					Message: "Must be true or false",
				}},
			})
			return
		}
	}

	if previewCountStr := c.Query(queryParamPreviewCount); previewCountStr != "" {
		if previewCount, err := strconv.Atoi(previewCountStr); err == nil {
			opts = append(opts, request.WithChatPreview(previewCount))
		} else {
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    "invalid_query_param",
				ErrorMessage: "Invalid query parameter",
				ErrorDetails: []restapi.ErrorDetail{{
					Field:   queryParamPreviewCount,
					Message: "Must be integer value",
				}},
			})
			return
		}
	}

	chats, err := h.service.GetByMemberID(c.Request.Context(), userID, opts...)
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, gin.H{
		"chats": chats,
	})
}

func (h *GenericChatHandler) GetChat(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	var opts []request.GetChatOption
	if enableLastUpdateID := c.Query(queryParamLastUpdateID); enableLastUpdateID != "" {
		switch enableLastUpdateID {
		case "true":
			opts = append(opts, request.WithChatLastUpdateID())
		case "false":
		default:
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    "invalid_query_param",
				ErrorMessage: "Invalid query parameter",
				ErrorDetails: []restapi.ErrorDetail{{
					Field:   queryParamLastUpdateID,
					Message: "Must be true or false",
				}},
			})
		}
	}

	if previewCountStr := c.Query(queryParamPreviewCount); previewCountStr != "" {
		if previewCount, err := strconv.Atoi(previewCountStr); err == nil {
			opts = append(opts, request.WithChatPreview(previewCount))
		} else {
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    "invalid_query_param",
				ErrorMessage: "Invalid query parameter",
				ErrorDetails: []restapi.ErrorDetail{{
					Field:   queryParamPreviewCount,
					Message: "Must be integer value",
				}},
			})
		}
	}

	chat, err := h.service.GetByChatID(c.Request.Context(), userID, chatId, opts...)
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, chat)
}
