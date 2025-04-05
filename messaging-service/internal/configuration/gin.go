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

	idemp := r.Group("/").
		Use(idempotency.New(
			idempotency.NewStorage(db.Redis, &idempotency.IdempotencyConfig{
				DataExp: 1 * time.Hour,
			}),
		))

	r.GET("/v1.0/chat/all", handlers.GenericChat.GetAllChats)
	r.GET("/v1.0/chat/:chatId", handlers.GenericChat.GetChat)

	idemp.POST("/v1.0/chat/personal", handlers.PersonalChat.CreateChat)
	r.PUT("/v1.0/chat/personal/:chatId/block", handlers.PersonalChat.BlockChat)
	r.PUT("/v1.0/chat/personal/:chatId/unblock", handlers.PersonalChat.UnblockChat)
	r.DELETE("/v1.0/chat/personal/:chatId", handlers.PersonalChat.DeleteChat)

	idemp.POST("/v1.0/chat/personal/secret", handlers.SecretPersonalChat.CreateChat)
	r.PUT("/v1.0/chat/personal/secret/:chatId/expiration", handlers.SecretGroup.SetExpiration)
	r.DELETE("/v1.0/chat/personal/secret/:chatId/:deleteMode", handlers.SecretPersonalChat.DeleteChat)

	idemp.POST("/v1.0/chat/group", handlers.GroupChat.CreateGroup)
	r.PUT("/v1.0/chat/group/:chatId", handlers.GroupChat.UpdateGroup)
	r.DELETE("/v1.0/chat/group/:chatId", handlers.GroupChat.DeleteGroup)
	r.PUT("/v1.0/chat/group/:chatId/member/:memberId", handlers.GroupChat.AddMember)
	r.DELETE("/v1.0/chat/group/:chatId/member/:memberId", handlers.GroupChat.DeleteMember)
	r.PUT("/v1.0/chat/group/:chatId/photo", handlers.GroupPhoto.UpdatePhoto)
	r.DELETE("/v1.0/chat/group/:chatId/photo", handlers.GroupPhoto.DeletePhoto)

	idemp.POST("/v1.0/chat/group/secret", handlers.SecretGroup.Create)
	r.PUT("/v1.0/chat/group/secret/:chatId", handlers.SecretGroup.Update)
	r.DELETE("/v1.0/chat/group/secret/:chatId", handlers.SecretGroup.Delete)
	r.PUT("/v1.0/chat/group/secret/:chatId/member/:memberId", handlers.SecretGroup.AddMember)
	r.DELETE("/v1.0/chat/group/secret/:chatId/member/:memberId", handlers.SecretGroup.DeleteMember)
	r.PUT("/v1.0/chat/group/secret/:chatId/photo", handlers.SecretGroupPhoto.UpdatePhoto)
	r.DELETE("/v1.0/chat/group/secret/:chatId/photo", handlers.SecretGroupPhoto.DeletePhoto)

	idemp.POST("/v1.0/chat/personal/:chatId/update/message/text", handlers.PersonalUpdate.SendTextMessage)
	r.DELETE("/v1.0/chat/personal/:chatId/update/message/:updateId/:deleteMode", handlers.PersonalUpdate.DeleteMessage)
	r.PUT("/v1.0/chat/personal/:chatId/update/message/text/:updateId", handlers.PersonalUpdate.EditTextMessage)
	idemp.POST("/v1.0/chat/personal/:chatId/update/message/file", handlers.PersonalFile.SendFileMessage)
	idemp.POST("/v1.0/chat/personal/:chatId/update/reaction", handlers.PersonalUpdate.SendReaction)
	r.DELETE("/v1.0/chat/personal/:chatId/update/reaction/:updateId", handlers.PersonalUpdate.DeleteReaction)
	idemp.POST("/v1.0/chat/personal/:chatId/update/text-message/forward", handlers.PersonalUpdate.ForwardTextMessage)
	idemp.POST("/v1.0/chat/personal/:chatId/update/file-message/forward", handlers.PersonalUpdate.ForwardFileMessage)

	idemp.POST("/v1.0/chat/group/:chatId/update/message/text", handlers.GroupUpdate.SendTextMessage)
	r.DELETE("/v1.0/chat/group/:chatId/update/message/:updateId/:deleteMode", handlers.GroupUpdate.DeleteMessage)
	r.PUT("/v1.0/chat/group/:chatId/update/message/text/:updateId", handlers.GroupUpdate.EditTextMessage)
	idemp.POST("/v1.0/chat/group/:chatId/update/message/file", handlers.GroupFile.SendFileMessage)
	idemp.POST("/v1.0/chat/group/:chatId/update/reaction", handlers.GroupUpdate.SendReaction)
	r.DELETE("/v1.0/chat/group/:chatId/update/reaction/:updateId", handlers.GroupUpdate.DeleteReaction)
	idemp.POST("/v1.0/chat/group/:chatId/update/text-message/forward", handlers.GroupUpdate.ForwardTextMessage)
	idemp.POST("/v1.0/chat/group/:chatId/update/file-message/forward", handlers.GroupUpdate.ForwardFileMessage)

	idemp.POST("/v1.0/chat/personal/secret/:chatId/update/secret", handlers.SecretPersonalUpdate.SendSecretUpdate)
	r.DELETE("/v1.0/chat/personal/secret/:chatId/update/secret/:updateId", handlers.SecretPersonalUpdate.DeleteSecretUpdate)

	idemp.POST("/v1.0/chat/group/secret/:chatId/update/secret", handlers.SecretGroupUpdate.SendSecretUpdate)
	r.DELETE("/v1.0/chat/group/secret/:chatId/update/secret/:updateId", handlers.SecretGroupUpdate.DeleteSecretUpdate)

	return r, nil
}
