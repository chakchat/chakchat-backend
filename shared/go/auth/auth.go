package auth

import (
	"context"
	"strings"

	"github.com/chakchat/chakchat/backend/shared/go/internal/restapi"
	"github.com/chakchat/chakchat/backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
)

const (
	headerAuthorization = "Authorization"
	keyClaims           = "_auth_claims"
)

// It is a copy of jwt claim names
const (
	ClaimName     = "name"
	ClaimUsername = "username"
	ClaimId       = "sub"
)

type Claims map[string]any

func NewJWT(conf *jwt.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get(headerAuthorization)
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok {
			restapi.SendUnautorized(c)
			c.Abort()
			return
		}

		claims, err := jwt.Parse(conf, jwt.Token(token))
		if err != nil {
			restapi.SendUnautorized(c)
			c.Abort()
			return
		}
		setClaims(c, Claims(claims))

		c.Next()
	}
}

func setClaims(c *gin.Context, claims Claims) {
	c.Set(keyClaims, claims)
}

func GetClaims(ctx context.Context) Claims {
	return ctx.Value(keyClaims).(Claims)
}
