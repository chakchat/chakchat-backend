package configuration

import (
	"net/http"
	"os"
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/rest/restapi"
	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/idempotency"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func GinEngine(handlers *Handlers, db *DB, conf *Config) (*gin.Engine, error) {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(otelgin.Middleware("messaging-service"))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeNotFound,
			ErrorMessage: "No such endpoint. Make sure that you use correct route and HTTP method.",
		})
	})

	jwtConf := &jwt.Config{
		SigningMethod: conf.JWT.SigningMethod,
		Issuer:        conf.JWT.Issuer,
		Audience:      conf.JWT.Audience,
		Type:          "internal_access",
	}
	jwtKey, err := os.ReadFile(conf.JWT.KeyFilePath)
	if err != nil {
		return nil, err
	}
	if err := jwtConf.RSAPublicOnlyKey(jwtKey); err != nil {
		return nil, err
	}
	r.Use(auth.NewJWT(&auth.JWTConfig{
		Conf:          jwtConf,
		DefaultHeader: "X-Internal-Token",
	}))

	idempotent := r.Group("/").
		Use(idempotency.New(
			idempotency.NewStorage(db.Redis, &idempotency.IdempotencyConfig{
				DataExp: 1 * time.Hour,
			}),
		))

	idempotent.POST("/v1.0/chat/personal", handlers.PersonalChat.CreateChat)
	r.PUT("/v1.0/chat/personal/:chatId/block", handlers.PersonalChat.BlockChat)
	r.PUT("/v1.0/chat/personal/:chatId/unblock", handlers.PersonalChat.UnblockChat)
	r.DELETE("/v1.0/chat/personal/:chatId", handlers.PersonalChat.DeleteChat)

	return r, nil
}
