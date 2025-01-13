package auth

import (
	"context"
	"log"
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

type JWTConfig struct {
	Conf          *jwt.Config
	Aud           string
	DefaultHeader string
}

func NewJWT(conf *JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var headerName string
		if conf.DefaultHeader != "" {
			headerName = conf.DefaultHeader
		} else {
			headerName = headerAuthorization
		}

		header := c.Request.Header.Get(headerName)
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok {
			log.Println("Unautorized because of invalid header")
			restapi.SendUnautorized(c)
			c.Abort()
			return
		}

		var claims jwt.Claims
		var err error
		if conf.Aud != "" {
			claims, err = jwt.ParseWithAud(conf.Conf, jwt.Token(token), conf.Aud)
		} else {
			claims, err = jwt.Parse(conf.Conf, jwt.Token(token))
		}
		if err != nil {
			log.Printf("Unauthorized: %s", err)
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
