package handlers

import (
	"context"
	"net/http"

	"github.com/chakchat/chakchat/backend/identity/internal/restapi"
	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/chakchat/chakchat/backend/identity/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type SignOutService interface {
	SignOut(ctx context.Context, refresh jwt.Token) error
}

func SignOut(service SignOutService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signOutRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		err := service.SignOut(c, jwt.Token(req.RefreshJWT))

		// I think that signing out expired token counts as a successful operation
		if err != nil && err != services.ErrRefreshTokenExpired {
			switch err {
			case services.ErrInvalidJWT:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidJWT, // Is it appropriate error_type?
					ErrorMessage: "Refresh token is invalid",
				})
			default:
				restapi.SendInternalError(c)
			}
			return
		}

		restapi.SendSuccess(c, signOutResponse{})
	}
}

type signOutRequest struct {
	RefreshJWT string `json:"refresh_token" binding:"required"`
}

type signOutResponse struct{}
