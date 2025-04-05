package update

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/gin-gonic/gin"
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

}
