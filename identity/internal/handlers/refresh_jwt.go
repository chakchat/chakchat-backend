package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/chakchat/chakchat/backend/identity/internal/restapi"
	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/chakchat/chakchat/backend/identity/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type RefreshJWTService interface {
	Refresh(ctx context.Context, refresh jwt.Token) (jwt.Pair, error)
}

func RefreshJWT(service RefreshJWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req refreshJWTRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		tokens, err := service.Refresh(c, jwt.Token(req.RefreshToken))

		if err != nil {
			log.Printf("met error in refresh-jwt: %s", err)
			switch err {
			case services.ErrRefreshTokenExpired:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeRefreshTokenExpired,
					ErrorMessage: "Refresh token expired",
				})
			case services.ErrRefreshTokenInvalidated:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeRefreshTokenInvalidated,
					ErrorMessage: "Refresh token invalidated",
				})
			case services.ErrInvalidTokenType:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidTokenType,
					ErrorMessage: "Invalid token type",
				})
			case services.ErrInvalidJWT:
				c.JSON(http.StatusBadRequest, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidJWT,
					ErrorMessage: "Invalid signature of JWT",
				})
			default:
				restapi.SendInternalError(c)
			}
			return
		}

		restapi.SendSuccess(c, refreshJWTResponse{
			AccessToken:  string(tokens.Access),
			RefreshToken: string(tokens.Refresh),
		})
	}
}

type refreshJWTRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshJWTResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
