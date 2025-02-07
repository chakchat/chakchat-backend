package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/chakchat/chakchat-backend/identity-service/internal/handlers"
	"github.com/chakchat/chakchat-backend/identity-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/identity-service/internal/services"
	"github.com/chakchat/chakchat-backend/identity-service/internal/sms"
	"github.com/chakchat/chakchat-backend/identity-service/internal/storage"
	"github.com/chakchat/chakchat-backend/identity-service/internal/userservice"
	"github.com/chakchat/chakchat-backend/shared/go/idempotency"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var conf *Config = loadConfig("/app/config.yml")

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %s", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown tracer provider: %s", err)
		}
	}()
	defer func() {
		if err := tp.ForceFlush(context.Background()); err != nil {
			log.Fatalf("ForceFlush failed: %s", err)
		}
	}()

	rdb := connectRedis()
	defer rdb.Close()
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Fatalf("Add instrument tracing to redis failed: %s", err)
	}

	sms := createSmsSender()
	usersClient, closeGrpc := createUsersClient()
	defer closeGrpc()

	accessTokenConfig := loadAccessTokenConfig()
	refreshTokenConfig := loadRefreshTokenConfig()
	internalTokenConfig := loadInternalTokenConfig()

	idempotencyStorage := createIdempotencyStorage(rdb)
	signInMetaStorage := createSignInMetaStorage(rdb)
	invalidatedTokenStorage := createInvalidatedTokenStorage(rdb)
	signUpMetaStorage := createSignUpMetaStorage(rdb)

	sendCodeService := createSignInSendCodeService(sms, signInMetaStorage, usersClient)
	signInService := services.NewSignInService(signInMetaStorage, accessTokenConfig, refreshTokenConfig)
	refreshService := services.NewRefreshService(invalidatedTokenStorage, accessTokenConfig, refreshTokenConfig)
	signOutService := services.NewSignOutService(invalidatedTokenStorage)
	identityService := services.NewIdentityService(accessTokenConfig, internalTokenConfig)
	signUpSendCodeService := createSignUpSendCodeService(sms, signUpMetaStorage, usersClient)
	signUpVerifyService := services.NewSignUpVerifyCodeService(signUpMetaStorage)
	signUpService := services.NewSignUpService(accessTokenConfig, refreshTokenConfig, usersClient, signUpMetaStorage)

	r := gin.New()

	r.Use(otelgin.Middleware("identity-service"))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeNotFound,
			ErrorMessage: "No such endpoint. Make sure that you use correct route and HTTP method.",
		})
	})

	r.Use(gin.Logger())

	r.Group("/v1.0").
		Use(idempotency.New(idempotencyStorage)).
		POST("/signin/send-phone-code", handlers.SignInSendCode(sendCodeService)).
		POST("/signin", handlers.SignIn(signInService)).
		POST("/refresh-token", handlers.RefreshJWT(refreshService)).
		POST("/signup/send-phone-code", handlers.SignUpSendCode(signUpSendCodeService)).
		POST("/signup/verify-code", handlers.SignUpVerifyCode(signUpVerifyService)).
		POST("/signup", handlers.SignUp(signUpService))

	r.PUT("/v1.0/sign-out", handlers.SignOut(signOutService))
	r.GET("/v1.0/identity", handlers.Identity(identityService))

	// Delete this line
	r.GET("/internal", func(c *gin.Context) {
		log.Println(c.Request.Header)
	})

	r.Run(":5000")
}

func createSignUpSendCodeService(sms services.SmsSender, storage *storage.SignUpMetaStorage,
	users userservice.UserServiceClient) *services.SignUpSendCodeService {
	config := &services.CodeConfig{
		SendFrequency: conf.PhoneCode.SendFrequency,
	}
	return services.NewSignUpSendCodeService(config, sms, storage, users)
}

func createSignUpMetaStorage(redisClient *redis.Client) *storage.SignUpMetaStorage {
	stConf := &storage.SignUpMetaConfig{
		MetaLifetime: conf.SignUpMeta.Lifetime,
	}
	return storage.NewSignUpMetaStorage(stConf, redisClient)
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

func createUsersClient() (client userservice.UserServiceClient, closeFunc func() error) {
	addr := conf.UserService.GrpcAddr
	// TODO: Insecure transport should be replaced in the future
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return userservice.NewUserServiceClient(conn), conn.Close
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

func createIdempotencyStorage(redisClient *redis.Client) idempotency.IdempotencyStorage {
	idempotencyConf := &idempotency.IdempotencyConfig{
		DataExp: conf.Idempotency.DataExp,
	}
	return idempotency.NewStorage(redisClient, idempotencyConf)
}

func createSignInSendCodeService(sms services.SmsSender, storage services.SignInMetaFindStorer,
	users userservice.UserServiceClient) *services.SignInSendCodeService {
	config := &services.CodeConfig{
		SendFrequency: conf.PhoneCode.SendFrequency,
	}
	return services.NewSignInSendCodeService(config, sms, storage, users)
}

func readKey(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(conf.Otlp.GrpcAddr),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("identity-service"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}
