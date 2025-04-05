package update

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/request"
	"github.com/gin-gonic/gin"
)

type SecretGroupUpdateService interface {
	SendSecretUpdate(context.Context, request.SendSecretUpdate) (*dto.SecretUpdateDTO, error)
	DeleteSecretUpdate(context.Context, request.DeleteSecretUpdate) (*dto.UpdateDeletedDTO, error)
}

type SecretGroupUpdateHandler struct {
	service SecretGroupUpdateService
}

func NewSecretGroupUpdateHandler(service SecretGroupUpdateService) *SecretGroupUpdateHandler {
	return &SecretGroupUpdateHandler{
		service: service,
	}
}

func (h *SecretGroupUpdateHandler) SendSecretUpdate(c *gin.Context) {

}
