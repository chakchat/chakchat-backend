package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/chakchat/chakchat-backend/user-service/internal/grpcservice"
	"github.com/chakchat/chakchat-backend/user-service/internal/handlers"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
}

func loadConfig(file string) *Config {
	viper.AutomaticEnv()

	viper.BindEnv("db.dsn", "DB_DSN")
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
	db, err := connectDB()
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	log.Println("connected to DB")

	userStorage := storage.NewUserStorage(db)
	userService := services.NewGetUserService(userStorage)
	userServer := handlers.NewUserServer(*userService)
	restrictionStorage := storage.NewRestrictionStorage(db)
	getUserService := services.NewGetService(userStorage, restrictionStorage)
	getRestrictionService := services.NewGetRestrictionService(restrictionStorage)
	updateUserService := services.NewUpdateUserService(userStorage)
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
		GET("/v1.0/me", getUserServer.GetMe()).
		GET("/v1.0/me/restrictions", handlers.GetRestrictions(getRestrictionService)).
		PUT("v1.0/me", handlers.UpdateUser(updateUserService, getUserService))
	r.GET("/v1.0/are-you-a-real-teapot", handlers.AmITeapot())

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

func connectDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(conf.DB.DSN))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.FieldRestriction{}, &models.FieldRestrictionUser{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}
	return db, nil
}

func readKey(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return key
}
