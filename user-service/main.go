package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/chakchat/chakchat-backend/user-service/internal/filestorage"
	"github.com/chakchat/chakchat-backend/user-service/internal/grpcservice"
	"github.com/chakchat/chakchat-backend/user-service/internal/handlers"
	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type JWTConfig struct {
	SigningMethod string        `mapstructure:"signing_method"`
	Lifetime      time.Duration `mapstructure:"lifetime"`
	Issuer        string        `mapstructure:"issuer"`
	Audience      []string      `mapstructure:"audience"`
	KeyFilePath   string        `mapstructure:"key_file_path"`
}

type Config struct {
	Jwt JWTConfig `mapstructure:"jwt"`

	DB struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"db"`

	Server struct {
		GRPCPort string `mapstructure:"grpc-port"`
	} `mapstructure:"server"`

	FileStorage struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"filestorage"`

	Otlp struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"otlp"`
}

func loadConfig(file string) *Config {
	viper.AutomaticEnv()

	viper.MustBindEnv("db.dsn", "DB_DSN")
	viper.BindEnv("server.port", "SERVER_PORT")

	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("config file not found: %v", err)
		} else {
			log.Fatalf("viper reading config failed: %s", err)
		}
	}

	config := new(Config)
	if err := viper.UnmarshalExact(&config); err != nil {
		log.Fatalf("viper config unmarshaling failed: %s", err)
	}

	return config
}

var conf *Config = loadConfig("/app/config.yml")

func main() {
	jwtConf := loadJWTConfig()

	db, err := pgxpool.New(context.Background(), conf.DB.DSN)
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	defer db.Close()
	log.Println("connected to DB")

	fileClient, close := createFileClient()
	defer close()

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

	userStorage := storage.NewUserStorage(db)
	userService := services.NewGetUserService(userStorage)
	userServer := handlers.NewUserServer(*userService)
	restrictionStorage := storage.NewRestrictionStorage(db)
	getUserService := services.NewGetService(userStorage, restrictionStorage)
	getRestrictionService := services.NewGetRestrictionService(restrictionStorage)
	updateUserService := services.NewUpdateUserService(userStorage)
	updateRestrictions := services.NewUpdateRestrService(restrictionStorage)
	processPhotoService := services.NewProcessPhotoService(userStorage, fileClient)
	getUserServer := handlers.NewGetUserHandler(getUserService)

	grpcPort := viper.GetString("server.grpc-port")

	listen, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcservice.RegisterUserServiceServer(grpcServer, userServer)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r := gin.New()

	r.Use(otelgin.Middleware("user-service"))

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
		Use(authMiddleware).
		GET("/v1.0/user/:userId", getUserServer.GetUserByID()).
		GET("/v1.0/user/username/:username", getUserServer.GetUserByUsername()).
		GET("/v1.0/users", getUserServer.GetUsersByCriteria()).
		GET("/v1.0/users/:users", getUserServer.GetUsers()).
		GET("/v1.0/me", getUserServer.GetMe()).
		GET("/v1.0/me/restrictions", handlers.GetRestrictions(getRestrictionService, getUserService)).
		PUT("v1.0/me", handlers.UpdateUser(updateUserService, getUserService)).
		PUT("v1.0/me/restrictions", handlers.UpdateRestrictions(updateRestrictions)).
		PUT("v1.0/me/profile-photo", handlers.UpdatePhoto(processPhotoService)).
		DELETE("v1.0/me/profile-photo", handlers.DeletePhoto(processPhotoService)).
		DELETE("v1.0/me", handlers.DeleteMe(updateUserService))
	r.Group("/").
		GET("/v1.0/are-you-a-real-teapot", handlers.AmITeapot()).
		GET("/v1.0/username/:username", handlers.CheckUserByUsername(getUserService))

	err = r.Run(":5004")
	if err != nil {
		log.Fatalf("Failed to run gin: %s", err)
	}
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

func readKey(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func createFileClient() (filestorage.FileStorageServiceClient, func() error) {
	addr := conf.FileStorage.GrpcAddr
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return filestorage.NewFileStorageServiceClient(conn), conn.Close
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
			semconv.ServiceNameKey.String("backend-services"),
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
