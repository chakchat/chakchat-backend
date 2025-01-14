package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chakchat/chakchat/backend/file-storage/internal/handlers"
	"github.com/chakchat/chakchat/backend/file-storage/internal/restapi"
	"github.com/chakchat/chakchat/backend/file-storage/internal/services"
	"github.com/chakchat/chakchat/backend/file-storage/internal/storage"
	"github.com/chakchat/chakchat/backend/shared/go/auth"
	"github.com/chakchat/chakchat/backend/shared/go/idempotency"
	"github.com/chakchat/chakchat/backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db := connectDB()
	redisClient := connectRedis()
	s3Client := connectS3()
	idempStorage := createIdempStorage(redisClient)
	fileMetaStorage := storage.NewFileMetaStorage(db)
	uploadMetaStorage := storage.NewUploadMetaStorage(db)

	uploadService := services.NewUploadService(fileMetaStorage, s3Client, s3Config)
	getFileService := services.NewGetFileService(fileMetaStorage)
	uploadInitService := services.NewUploadInitService(uploadMetaStorage, s3Client, s3Config)
	uploadPartService := services.NewUploadPartService(uploadMetaStorage, s3Client, s3Config)
	uploadAbortService := services.NewUploadAbortService(uploadMetaStorage, s3Client, s3Config)

	r := gin.New()

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
		POST("/v1.0")

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
