package handler

import (
	"errors"
	"net/http"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/services"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/storage"
	"github.com/gin-gonic/gin"
)

type OnlineStatusServer struct {
	service *services.StatusService
}

func NewOnlineStatusServer(service *services.StatusService) *OnlineStatusServer {
	return &OnlineStatusServer{
		service: service,
	}
}

func (s *OnlineStatusServer) GetStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIds := c.QueryArray("users")
		if len(userIds) == 0 {
			c.JSON(http.StatusBadRequest, restapi.ErrTypeBadRequest)
			return
		}

		status, err := s.service.GetStatus(c.Request.Context(), userIds)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				c.JSON(http.StatusNotFound, restapi.ErrTypeNotFound)
				return
			}
		}

		restapi.SendSuccess(c, status)
	}
}
