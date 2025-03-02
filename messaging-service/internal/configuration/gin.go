package configuration

import (
	"log"
	"os"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/chakchat/chakchat/backend/shared/go/idempotency"
	"github.com/gin-gonic/gin"
)

func GinEngine(handlers *Handlers, db *DB, conf *Config) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(auth.NewJWT(&auth.JWTConfig{
		Conf: &jwt.Config{
			SigningMethod: conf.JWT.SigningMethod,
			Issuer:        conf.JWT.Issuer,
			Audience:      conf.JWT.Audience,
			Type:          "internal_access",
			SymmetricKey:  readKey(conf.JWT.KeyFilePath),
		},
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
	r.DELETE("/v1.0/chat/personal/:chatId/delete/{:deleteMode}", handlers.PersonalChat.DeleteChat)

	return r
}

func readKey(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot read file: %s", err)
	}
	return key
}
