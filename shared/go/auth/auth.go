package auth

import (
	"context"
	"log"
	"strings"

	"github.com/chakchat/chakchat-backend/shared/go/internal/restapi"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
)

const (
	headerAuthorization = "Authorization"
)

// It is a copy of jwt claim names
const (
	ClaimName     = "name"
	ClaimUsername = "username"
	ClaimId       = "sub"
)

type keyClaimsType int

var keyClaims = keyClaimsType(69)

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

		ctx := context.WithValue(
			c.Request.Context(),
			keyClaims, Claims(claims),
		)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetClaims(ctx context.Context) Claims {
	if val := ctx.Value(keyClaims); val != nil {
		return val.(Claims)
	}
	return nil
}
