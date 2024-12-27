package handlers

import (
	"net/http"
	"strings"

	"github.com/chakchat/chakchat/backend/identity/internal/jwt"
	"github.com/chakchat/chakchat/backend/identity/internal/restapi"
	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	HeaderAuthorization = "Authorization"
	HeaderInternalToken = "X-Internal-Token"
)

type IdentityService interface {
	Idenitfy(access jwt.JWT) (jwt.InternalJWT, error)
}

func GetIdentity(service IdentityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(HeaderAuthorization)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Authorization header must contain authorization info",
			})
			return
		}

		publicToken, ok := extractJWT(authHeader)
		if !ok {
			c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
				ErrorType:    restapi.ErrTypeUnautorized,
				ErrorMessage: "Invalid Authorization header format",
			})
		}

		internalToken, err := service.Idenitfy(publicToken)
		if err != nil {
			switch err {
			case services.ErrInvalidJWT:
				c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUnautorized,
					ErrorMessage: "Invalid Authorization token",
				})
			case services.ErrAccessTokenExpired:
				c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUnautorized,
					ErrorMessage: "Access token expired",
				})
			default:
				restapi.SendInternalError(c)
			}
			return
		}

		c.Header(HeaderInternalToken, "Bearer "+string(internalToken))
		c.Writer.WriteHeader(http.StatusNoContent)
	}
}

func extractJWT(authHeader string) (jwt.JWT, bool) {
	found, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		return "", false
	}
	return jwt.JWT(found), true
}