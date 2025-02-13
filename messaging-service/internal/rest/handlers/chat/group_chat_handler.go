package chat

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const paramMemberID = "memberId"

type GroupChatService interface {
	CreateGroup(ctx context.Context, req request.CreateGroup) (*dto.GroupChatDTO, error)
	UpdateGroupInfo(ctx context.Context, req request.UpdateGroupInfo) (*dto.GroupChatDTO, error)
	DeleteGroup(ctx context.Context, req request.DeleteChat) error
	AddMember(ctx context.Context, req request.AddMember) (*dto.GroupChatDTO, error)
	DeleteMember(ctx context.Context, req request.DeleteMember) (*dto.GroupChatDTO, error)
}

type GroupChatHandler struct {
	service GroupChatService
}

func NewGroupChatHandler(service GroupChatService) *GroupChatHandler {
	return &GroupChatHandler{
		service: service,
	}
}

func (h *GroupChatHandler) CreateGroup(c *gin.Context) {
	userId := getUserID(c.Request.Context())

	req := struct {
		Name    string      `json:"name"`
		Members []uuid.UUID `json:"members"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	group, err := h.service.CreateGroup(c.Request.Context(), request.CreateGroup{
		SenderID: userId,
		Members:  req.Members,
		Name:     req.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	restapi.SendSuccess(c, newGroupResponse(group))
}

func (h *GroupChatHandler) UpdateGroup(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	req := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	group, err := h.service.UpdateGroupInfo(c.Request.Context(), request.UpdateGroupInfo{
		ChatID:      chatId,
		SenderID:    userId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.Error(err)
		return
	}

	restapi.SendSuccess(c, newGroupResponse(group))
}

func (h GroupChatHandler) DeleteGroup(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	err = h.service.DeleteGroup(c.Request.Context(), request.DeleteChat{
		ChatID:   chatId,
		SenderID: userId,
	})
	if err != nil {
		c.Error(err)
		return
	}

	restapi.SendSuccess(c, struct{}{})
}

func (h *GroupChatHandler) AddMember(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	memberId, err := uuid.Parse(c.Param(paramMemberID))
	if err != nil {
		restapi.SendInvalidMemberID(c)
		return
	}

	group, err := h.service.AddMember(c.Request.Context(), request.AddMember{
		ChatID:   chatId,
		SenderID: userId,
		MemberID: memberId,
	})
	if err != nil {
		c.Error(err)
		return
	}

	restapi.SendSuccess(c, newGroupResponse(group))
}

func (h *GroupChatHandler) DeleteMember(c *gin.Context) {
	chatId, err := uuid.Parse(c.Param(paramChatID))
	if err != nil {
		restapi.SendInvalidChatID(c)
		return
	}
	userId := getUserID(c.Request.Context())

	memberId, err := uuid.Parse(c.Param(paramMemberID))
	if err != nil {
		restapi.SendInvalidMemberID(c)
		return
	}

	group, err := h.service.DeleteMember(c.Request.Context(), request.DeleteMember{
		ChatID:   chatId,
		SenderID: userId,
		MemberID: memberId,
	})
	if err != nil {
		c.Error(err)
		return
	}

	restapi.SendSuccess(c, newGroupResponse(group))
}

type groupResponse struct {
	ChatID  uuid.UUID   `json:"chat_id"`
	Admin   uuid.UUID   `json:"admin"`
	Members []uuid.UUID `json:"members"`

	Name        string `json:"name"`
	Description string `json:"description"`
	GroupPhoto  string `json:"group_photo"`
	CreatedAt   int64  `json:"created_at"`
}

func newGroupResponse(dto *dto.GroupChatDTO) groupResponse {
	return groupResponse{
		ChatID:      dto.ID,
		Admin:       dto.Admin,
		Members:     dto.Members,
		Name:        dto.Name,
		Description: dto.Description,
		GroupPhoto:  dto.GroupPhoto,
		CreatedAt:   dto.CreatedAt,
	}
}
