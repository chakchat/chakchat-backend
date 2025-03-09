package chat

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenericChatService interface {
	GetByMemberID(ctx context.Context, memberID uuid.UUID) ([]services.GenericChat, error)
	GetByChatID(ctx context.Context, senderID, chatID uuid.UUID) (*services.GenericChat, error)
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

	chats, err := h.service.GetByMemberID(c.Request.Context(), userID)
	if err != nil {
		resp := errmap.Map(err)
		c.JSON(resp.Code, resp.Body)
		return
	}

	resp := struct {
		Chats []any `json:"chats"`
	}{
		Chats: make([]any, 0, len(chats)),
	}

	for _, chat := range chats {
		resp.Chats = append(resp.Chats, convertGenericChatResp(&chat))
	}

	restapi.SendSuccess(c, resp)
}

func (h *GenericChatHandler) GetChat(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userID := getUserID(c.Request.Context())

	chat, err := h.service.GetByChatID(c.Request.Context(), userID, chatId)
	if err != nil {
		resp := errmap.Map(err)
		c.JSON(resp.Code, resp.Body)
		return
	}

	restapi.SendSuccess(c, convertGenericChatResp(chat))
}
