package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/handlers"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/services"
	"github.com/chakchat/chakchat-backend/file-storage-service/internal/storage"
	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/idempotency"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var conf *Config = loadConfig("/app/config.yml")

func main() {
	jwtConf := loadJWTConfig()
	uploadConfig := &handlers.UploadConfig{
		FileSizeLimit: conf.Upload.FileSizeLimit,
	}
	multipartConfig := &handlers.MultipartUploadConfig{
		MinFileSize: conf.MultipartUpload.MinFileSize,
		MaxPartSize: conf.MultipartUpload.MaxPartSize,
	}
	s3Config := &services.S3Config{
		Bucket:    conf.S3.Bucket,
		UrlPrefix: conf.S3.UrlPrefix,
	}

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

	db := connectDB()
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Fatalf("Add tracing to gorm failed: %s", err)
	}

	s3Client := connectS3()
	idempStorage := createIdempStorage(rdb)
	fileMetaStorage := storage.NewFileMetaStorage(db)
	uploadMetaStorage := storage.NewUploadMetaStorage(db)

	uploadService := services.NewUploadService(fileMetaStorage, s3Client, s3Config)
	getFileService := services.NewGetFileService(fileMetaStorage)
	uploadInitService := services.NewUploadInitService(uploadMetaStorage, s3Client, s3Config)
	uploadPartService := services.NewUploadPartService(uploadMetaStorage, s3Client, s3Config)
	uploadAbortService := services.NewUploadAbortService(uploadMetaStorage, s3Client, s3Config)
	uploadCompleteService := services.NewUploadCompleteService(fileMetaStorage, uploadMetaStorage, s3Client, s3Config)

	r := gin.New()

	r.Use(otelgin.Middleware("file-storage-service"))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, restapi.ErrorResponse{
			ErrorType:    restapi.ErrTypeNotFound,
			ErrorMessage: "No such endpoint. Make sure that you use correct route and HTTP method.",
		})
	})

	authMiddleware := auth.NewJWT(&auth.JWTConfig{
		Conf:          jwtConf,
		DefaultHeader: "X-Internal-Token",
	})

	r.Group("/").
		Use(idempotency.New(idempStorage)).
		Use(authMiddleware).
		POST("/v1.0/upload", handlers.Upload(uploadConfig, uploadService)).
		POST("/v1.0/upload/multipart/init", handlers.UploadInit(multipartConfig, uploadInitService)).
		POST("/v1.0/upload/multipart/complete", handlers.UploadComplete(uploadCompleteService))

	r.Group("/").
		Use(authMiddleware).
		GET("/v1.0/file/:fileId", handlers.GetFile(getFileService)).
		PUT("/v1.0/upload/multipart/part", handlers.UploadPart(multipartConfig, uploadPartService)).
		PUT("/v1.0/upload/multipart/abort", handlers.UploadAbort(uploadAbortService))

	r.Run(":5004")
}

func loadJWTConfig() *jwt.Config {
	config := &jwt.Config{
		SigningMethod: conf.Jwt.SigningMethod,
		Lifetime:      conf.Jwt.Lifetime,
		Issuer:        conf.Jwt.Issuer,
		Audience:      conf.Jwt.Audience,
		Type:          "internal_access",
	}
	if err := config.RSAPublicOnlyKey(readKey(conf.Jwt.KeyFilePath)); err != nil {
		log.Fatal(err)
	}
	return config
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(conf.DB.DSN))
	if err != nil {
		log.Fatalf("database connecting failed: %s", err)
	}
	db.AutoMigrate(&storage.FileMeta{})
	db.AutoMigrate(&storage.UploadMeta{})
	return db
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

func connectS3() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			conf.AWS.AccessKeyId,
			conf.AWS.SecretAccessKey,
			"")),
		config.WithBaseEndpoint(conf.AWS.EndpointUrl),
	)
	if err != nil {
		log.Fatalf("Err loading default config: %s", err)
	}
	return s3.NewFromConfig(cfg)
}

func createIdempStorage(redisClient *redis.Client) idempotency.IdempotencyStorage {
	return idempotency.NewStorage(redisClient, &idempotency.IdempotencyConfig{
		DataExp: conf.Idempotency.DataExp,
	})
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
			semconv.ServiceNameKey.String("backend-service"),
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
