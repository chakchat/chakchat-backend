package handlers

import (
	"context"
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
	Idenitfy(ctx context.Context, access jwt.Token) (jwt.InternalToken, error)
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

		internalToken, err := service.Idenitfy(c, publicToken)
		if err != nil {
			switch err {
			case services.ErrInvalidJWT:
				c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeUnautorized,
					ErrorMessage: "Invalid Authorization token",
				})
			case services.ErrAccessTokenExpired:
				c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeAccessTokenExpired,
					ErrorMessage: "Access token expired",
				})
			case services.ErrInvalidTokenType:
				c.JSON(http.StatusUnauthorized, restapi.ErrorResponse{
					ErrorType:    restapi.ErrTypeInvalidTokenType,
					ErrorMessage: "Invalid token type",
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

func extractJWT(authHeader string) (jwt.Token, bool) {
	found, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		return "", false
	}
	return jwt.Token(found), true
}
