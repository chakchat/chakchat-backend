package main

import (
	"fmt"
	"log"
	"net"

	"github.com/chakchat/chakchat-backend/user-service/internal/grpcservice"
	"github.com/chakchat/chakchat-backend/user-service/internal/handlers"
	"github.com/chakchat/chakchat-backend/user-service/internal/models"
	"github.com/chakchat/chakchat-backend/user-service/internal/services"
	"github.com/chakchat/chakchat-backend/user-service/internal/storage"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
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

	db, err := connectDB()
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	userStorage := storage.NewUserStorage(db)
	userService := services.NewGetUserService(userStorage)
	userServer := handlers.NewUserServer(*userService)

	grpcPort := viper.GetString("server.grpc-port")

	listen, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcservice.RegisterUserServiceServer(grpcServer, userServer)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func connectDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(conf.DB.DSN))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}
	return db, nil
}
