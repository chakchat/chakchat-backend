package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/domain/services"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/infrastructure/messages"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/infrastructure/ws"
	"github.com/chakchat/chakchat-backend/live-connection-service/restapi"
	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type JWTConfig struct {
	SigningMethod string        `mapstructure:"signing_method"`
	Lifetime      time.Duration `mapstructure:"lifetime"`
	Issuer        string        `mapstructure:"issuer"`
	Audience      []string      `mapstructure:"audience"`
	KeyFilePath   string        `mapstructure:"key_file_path"`
}

type Config struct {
	Kafka struct {
		Brokers      []string `mapstructure:"btokers"`
		ConsumeTopic string   `mapstructure:"consume_topic"`
		ProduceTopic string   `mapstructure:"produce_topic"`
	} `mapstructure:"kafka"`

	Jwt JWTConfig `mapstructure:"jwt"`
}

func loadConfig(file string) *Config {
	viper.AutomaticEnv()

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

var conf = loadConfig("/app/config.yml")

func main() {

	jwtConf := loadJWTConfig()

	hub := ws.NewHub()

	kafkaProducer := messages.NewProducer(messages.ProducerConfig{
		Brokers: conf.Kafka.Brokers,
		Topic:   conf.Kafka.ProduceTopic,
	})
	defer kafkaProducer.Close()

	kafkaConsumer := messages.NewConsumer(&messages.ConsumerConf{
		Brokers: conf.Kafka.Brokers,
		Topic:   conf.Kafka.ConsumeTopic,
		GroupID: "live-connection-group",
	})

	defer kafkaConsumer.Stop()

	messageProcessor := services.NewKafkaProcessor(hub, kafkaProducer)

	go kafkaConsumer.Start(context.Background(), messageProcessor.MessageHandler)

	r := gin.New()
	r.Use(otelgin.Middleware("live-connection-service"))

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
		GET("/ws", hub.WebSocketHandler()).
		GET("/health", hub.HealthCheck())

	err := r.Run(":5004")
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
