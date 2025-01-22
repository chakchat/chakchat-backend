package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const paramChatId = "chatId"

type InChatNotBlockedChecker interface {
	CheckChatNotBlocked(userId, chatId uuid.UUID) error
}

type InChatChecker interface {
	CheckChat(userId, chatId uuid.UUID) error
}

type GroupAdminChecker interface {
	CheckAdmin(userId, chatId uuid.UUID) error
}

// Checks if user is in chat and not blocked
func NotBlockedAuthorization(ch InChatNotBlockedChecker) gin.HandlerFunc {
	return chatActionAuthorization(ch.CheckChatNotBlocked)
}

func InChatAuthorization(ch InChatChecker) gin.HandlerFunc {
	return chatActionAuthorization(ch.CheckChat)
}

func GroupAdminAuthorization(ch GroupAdminChecker) gin.HandlerFunc {
	return chatActionAuthorization(ch.CheckAdmin)
}

func chatActionAuthorization(checker func(userId, chatId uuid.UUID) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Abort()

		userId, ok := getUserId(c)
		if !ok {
			// This may happen if userId claim is not even set.
			// Maybe authentication middleware didn't handle this request.
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Unautorized. No userId claim",
			})
			return
		}

		chatIdParam := c.Param(paramChatId)
		if chatIdParam == "" {
			log.Printf("%s param is not found. But all these endpoints must require this param", paramChatId)
			// I think it is internal server error because this middleware
			// because this middleware must not be executed on endpoints that doesn't require this authorization policy
			restapi.SendInternalError(c)
			return
		}
		chatId, err := uuid.Parse(chatIdParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeInvalidParam,
				ErrorMessage: "Invalid " + paramChatId + " uuid parameter",
			})
			return
		}

		if err := checker(userId, chatId); err != nil {
			c.JSON(http.StatusForbidden, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeForbidden,
				ErrorMessage: "Forbidden. You doesn't have permission to perform this action",
			})
			return
		}

		c.Next()
	}
}

func getUserId(ctx context.Context) (uuid.UUID, bool) {
	claims := auth.GetClaims(ctx)

	idClaim, ok := claims[auth.ClaimId]
	if !ok {
		return uuid.Nil, false
	}

	userId, err := uuid.Parse(idClaim.(string))
	if err != nil {
		log.Printf("userId claim parsing failed: %s", err)
		return uuid.Nil, false
	}

	return userId, true
}
