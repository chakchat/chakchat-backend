package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/chakchat/chakchat/backend/identity/internal/handlers"
	"github.com/chakchat/chakchat/backend/identity/internal/idempotency"
	"github.com/chakchat/chakchat/backend/identity/internal/jwt"
	"github.com/chakchat/chakchat/backend/identity/internal/restapi"
	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/chakchat/chakchat/backend/identity/internal/sms"
	"github.com/chakchat/chakchat/backend/identity/internal/storage"
	"github.com/chakchat/chakchat/backend/identity/internal/userservice"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conf *Config = loadConfig("/app/config.yml")

func main() {
	r := gin.New()

	redisClient := connectRedis()
	sms := createSmsSender()
	usersClient := createUsersClient()

	accessTokenConfig := loadAccessTokenConfig()
	refreshTokenConfig := loadRefreshTokenConfig()
	internalTokenConfig := loadInternalTokenConfig()

	idempotencyStorage := createIdempotencyStorage(redisClient)
	signInMetaStorage := createSignInMetaStorage(redisClient)
	invalidatedTokenStorage := createInvalidatedTokenStorage(redisClient)

	sendCodeService := createSendCodeService(sms, signInMetaStorage, usersClient)
	signInService := services.NewSignInService(signInMetaStorage, accessTokenConfig, refreshTokenConfig)
	refreshService := services.NewRefreshService(invalidatedTokenStorage, accessTokenConfig, refreshTokenConfig)
	signOutService := services.NewSignOutService(invalidatedTokenStorage)
	identityService := services.NewIdentityService(accessTokenConfig, internalTokenConfig)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeNotFound,
			ErrorMessage: "No such endpoint. Make sure that you use correct route and HTTP method.",
		})
	})

	r.Group("/v1.0").
		Use(idempotency.New(idempotencyStorage)).
		POST("/signin/send-phone-code", handlers.SendCode(sendCodeService)).
		POST("/signin", handlers.SignIn(signInService)).
		POST("/refresh-token", handlers.RefreshJWT(refreshService))

	r.Use(gin.Logger())
	r.PUT("/v1.0/sign-out", handlers.SignOut(signOutService))
	r.GET("/v1.0/identity", handlers.Identity(identityService))

	// Delete this line
	r.GET("/internal", func(c *gin.Context) {
		log.Println(c.Request.Header)
	})

	r.Run(":5000")
}

func createSmsSender() services.SmsSender {
	if conf.Sms.Type == "stub" {
		return sms.NewSmsServerStubSender(conf.Sms.Stub.Addr)
	}
	return &sms.SmsSenderFake{}
}

func createInvalidatedTokenStorage(redisClient *redis.Client) *storage.InvalidatedTokenStorage {
	conf := &storage.InvalidatedTokenConfig{
		InvalidatedExp: conf.InvalidatedTokenStorage.Exp,
	}
	return storage.NewInvalidatedTokenStorage(conf, redisClient)
}

func loadAccessTokenConfig() *jwt.Config {
	return &jwt.Config{
		SigningMethod: conf.AccessToken.SigningMethod,
		Lifetime:      conf.AccessToken.Lifetime,
		Issuer:        conf.AccessToken.Issuer,
		Audience:      conf.AccessToken.Audience,
		Type:          "access",
		SymmetricKey:  readKey(conf.AccessToken.KeyFilePath),
	}
}

func loadRefreshTokenConfig() *jwt.Config {
	return &jwt.Config{
		SigningMethod: conf.RefreshToken.SigningMethod,
		Lifetime:      conf.RefreshToken.Lifetime,
		Issuer:        conf.RefreshToken.Issuer,
		Audience:      conf.RefreshToken.Audience,
		Type:          "refresh",
		SymmetricKey:  readKey(conf.RefreshToken.KeyFilePath),
	}
}

func loadInternalTokenConfig() *jwt.Config {
	res := &jwt.Config{
		SigningMethod: conf.InternalToken.SigningMethod,
		Lifetime:      conf.InternalToken.Lifetime,
		Issuer:        conf.InternalToken.Issuer,
		Audience:      conf.InternalToken.Audience,
		Type:          "internal_access",
	}
	res.RSAKeys(readKey(conf.InternalToken.KeyFilePath))
	return res
}

func createUsersClient() userservice.UserServiceClient {
	addr := conf.UserService.GrpcAddr
	// TODO: Insecure transport should be replaced in the future
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return userservice.NewUserServiceClient(conn)
}

func connectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis connection establishing failed: %s", err)
	}
	log.Println("redis connection established")
	return client
}

func createSignInMetaStorage(redisClient *redis.Client) *storage.SignInMetaStorage {
	config := &storage.SignInMetaConfig{
		MetaLifetime: conf.SignInMeta.Lifetime,
	}
	return storage.NewSignInMetaStorage(config, redisClient)
}

func createIdempotencyStorage(redisClient *redis.Client) *storage.IdempotencyStorage {
	idempotencyConf := &storage.IdempotencyConfig{
		DataExp: conf.Idempotency.DataExp,
	}
	return storage.NewIdempotencyStorage(redisClient, idempotencyConf)
}

func createSendCodeService(sms services.SmsSender, storage services.MetaFindStorer,
	users userservice.UserServiceClient) *services.SendCodeService {
	config := &services.CodeConfig{
		SendFrequency: conf.PhoneCode.SendFrequency,
	}
	return services.NewSendCodeService(config, sms, storage, users)
}

func readKey(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return key
}
