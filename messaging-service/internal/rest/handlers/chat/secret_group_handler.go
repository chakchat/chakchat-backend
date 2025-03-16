package chat

import (
	"context"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/errmap"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/response"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecretGroupService interface {
	CreateGroup(ctx context.Context, req request.CreateSecretGroup) (*dto.SecretGroupChatDTO, error)
	UpdateGroupInfo(ctx context.Context, req request.UpdateSecretGroupInfo) (*dto.SecretGroupChatDTO, error)
	DeleteGroup(ctx context.Context, req request.DeleteChat) error
	AddMember(ctx context.Context, req request.AddMember) (*dto.SecretGroupChatDTO, error)
	DeleteMember(ctx context.Context, req request.DeleteMember) (*dto.SecretGroupChatDTO, error)
	SetExpiration(ctx context.Context, req request.SetExpiration) (*dto.SecretGroupChatDTO, error)
}

type SecretGroupHandler struct {
	service SecretGroupService
}

func NewSecretGroupHandler(service SecretGroupService) *SecretGroupHandler {
	return &SecretGroupHandler{
		service: service,
	}
}

func (h *SecretGroupHandler) Create(c *gin.Context) {
	userId := getUserID(c.Request.Context())

	req := struct {
		Name    string      `json:"name"`
		Members []uuid.UUID `json:"members"`
	}{}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		restapi.SendUnprocessableJSON(c)
		return
	}

	group, err := h.service.CreateGroup(c.Request.Context(), request.CreateSecretGroup{
		SenderID: userId,
		Members:  req.Members,
		Name:     req.Name,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.SecretGroupChat(group))
}

func (h *SecretGroupHandler) Update(c *gin.Context) {
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

	group, err := h.service.UpdateGroupInfo(c.Request.Context(), request.UpdateSecretGroupInfo{
		ChatID:      chatId,
		SenderID:    userId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.SecretGroupChat(group))
}

func (h *SecretGroupHandler) AddMember(c *gin.Context) {
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
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.SecretGroupChat(group))
}

func (h *SecretGroupHandler) DeleteMember(c *gin.Context) {
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
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.SecretGroupChat(group))
}

func (h *SecretGroupHandler) Delete(c *gin.Context) {
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
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, struct{}{})
}

func (h *SecretGroupHandler) SetExpiration(c *gin.Context) {
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

	group, err := h.service.SetExpiration(c.Request.Context(), request.SetExpiration{
		ChatID:     chatId,
		SenderID:   userId,
		Expiration: req.Expiration,
	})
	if err != nil {
		errmap.Respond(c, err)
		return
	}

	restapi.SendSuccess(c, response.SecretGroupChat(group))
}

// type secretGroupResponse struct {
// 	ChatID    uuid.UUID `json:"chat_id"`
// 	CreatedAt int64     `json:"created_at"`

// 	AdminID uuid.UUID   `json:"admin_id"`
// 	Members []uuid.UUID `json:"members"`

// 	Name          string `json:"name"`
// 	Description   string `json:"description"`
// 	GroupPhotoURL string `json:"group_photo_url"`

// 	Expiration *time.Duration `json:"expiration"`
// }

// func newSecretGroupResponse(dto *dto.SecretGroupChatDTO) secretGroupResponse {
// 	return secretGroupResponse{
// 		ChatID:        dto.ID,
// 		CreatedAt:     dto.CreatedAt,
// 		AdminID:       dto.Admin,
// 		Members:       dto.Members,
// 		Name:          dto.Name,
// 		Description:   dto.Description,
// 		GroupPhotoURL: dto.GroupPhotoURL,
// 		Expiration:    dto.Expiration,
// 	}
// }
