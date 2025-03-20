package handlers

import (
	"context"
	"log"

	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/gin-gonic/gin"
)

type CheckUserServer interface {
	CheckUserByUsername(ctx context.Context, username string) (*bool, error)
}

func CheckUserByUsername(service CheckUserServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Start checking")
		username := c.Param("username")
		log.Println("Got parametr")

		res, err := service.CheckUserByUsername(c.Request.Context(), username)
		if err != nil {
			restapi.SendInternalError(c)
			return
		}
		log.Println("Return")

		restapi.SendSuccess(c, gin.H{"user_exists": res})
	}
}
