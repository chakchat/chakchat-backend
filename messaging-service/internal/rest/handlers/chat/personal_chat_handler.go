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

const paramChatID = "chatId"

type PersonalChatService interface {
	BlockChat(ctx context.Context, req request.BlockChat) (*dto.PersonalChatDTO, error)
	UnblockChat(ctx context.Context, req request.UnblockChat) (*dto.PersonalChatDTO, error)
	CreateChat(ctx context.Context, req request.CreatePersonalChat) (*dto.PersonalChatDTO, error)
	DeleteChat(ctx context.Context, req request.DeleteChat) error
}

type PersonalChatHandler struct {
	service PersonalChatService
}

func NewPersonalChatHandler(service PersonalChatService) *PersonalChatHandler {
	return &PersonalChatHandler{
		service: service,
	}
}

func (h *PersonalChatHandler) CreateChat(c *gin.Context) {
	userId := getUserID(c.Request.Context())

	req := struct {
		MemberID uuid.UUID `json:"member_id"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	chat, err := h.service.CreateChat(c.Request.Context(), request.CreatePersonalChat{
		SenderID: userId,
		MemberID: req.MemberID,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.PersonalChat(chat))
}

func (h *PersonalChatHandler) BlockChat(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	chat, err := h.service.BlockChat(c.Request.Context(), request.BlockChat{
		ChatID:   chatId,
		SenderID: userId,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.PersonalChat(chat))
}

func (h *PersonalChatHandler) UnblockChat(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	chat, err := h.service.UnblockChat(c.Request.Context(), request.UnblockChat{
		ChatID:   chatId,
		SenderID: userId,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.PersonalChat(chat))
}

func (h *PersonalChatHandler) DeleteChat(c *gin.Context) {
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

// type personalChatResponse struct {
// 	ID      uuid.UUID    `json:"chat_id"`
// 	Members [2]uuid.UUID `json:"members"`

// 	Blocked   bool        `json:"blocked"`
// 	BlockedBy []uuid.UUID `json:"blocked_by"`
// 	CreatedAt int64       `json:"created_at"`
// }

// func newPersonalChatResponse(dto *dto.PersonalChatDTO) personalChatResponse {
// 	return personalChatResponse{
// 		ID:        dto.ID,
// 		Members:   dto.Members,
// 		Blocked:   dto.Blocked,
// 		BlockedBy: dto.BlockedBy,
// 		CreatedAt: dto.CreatedAt,
// 	}
// }
