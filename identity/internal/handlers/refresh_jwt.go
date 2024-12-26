package handlers

import (
	"net/http"

	"github.com/chakchat/chakchat/backend/identity/internal/jwt"
	"github.com/chakchat/chakchat/backend/identity/internal/restapi"
	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/gin-gonic/gin"
)

type RefreshJWTService interface {
	Refresh(refresh jwt.JWT) (jwt.Pair, error)
}

func RefreshJWT(service RefreshJWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req refreshJWTRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			restapi.SendUnprocessableJSON(c)
			return
		}

		tokens, err := service.Refresh(jwt.JWT(req.RefreshToken))

		if err != nil {
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
