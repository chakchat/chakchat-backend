package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/handler"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/mq"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/restapi"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/services"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/storage"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/ws"
	"github.com/chakchat/chakchat-backend/shared/go/auth"
	"github.com/chakchat/chakchat-backend/shared/go/jwt"
	"github.com/chakchat/chakchat-backend/shared/go/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type JWTConfig struct {
	SigningMethod string        `mapstructure:"signing_method"`
	Lifetime      time.Duration `mapstructure:"lifetime"`
	Issuer        string        `mapstructure:"issuer"`
	Audience      []string      `mapstructure:"audience"`
	KeyFilePath   string        `mapstructure:"key_file_path"`
}

type Config struct {
	ConsumeKafka struct {
		Brokers      []string `mapstructure:"brokers"`
		Topic string   `mapstructure:"topic"`
	} `mapstructure:"consume_kafka"`

	ProduceKafka struct {
		Brokers      []string `mapstructure:"brokers"`
		Topic string   `mapstructure:"topic"`
	} `mapstructure:"produce_kafka"`

	Jwt JWTConfig `mapstructure:"jwt"`

	DB struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"db"`
}

func loadConfig(file string) *Config {
	viper.AutomaticEnv()

	viper.MustBindEnv("db.dsn", "PG_CONN_STRING")

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

	pgxDb, err := pgxpool.New(context.Background(), conf.DB.DSN)
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	defer pgxDb.Close()
	db := postgres.Tracing(pgxDb)
	log.Println("connected to DB")

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

	hub := ws.NewHub()

	kafkaProducer := mq.NewProducer(mq.ProducerConfig{
		Brokers: conf.ProduceKafka.Brokers,
		Topic:   conf.ProduceKafka.Topic,
	})
	defer kafkaProducer.Close()

	kafkaConsumer := mq.NewConsumer(&mq.ConsumerConf{
		Brokers: conf.ConsumeKafka.Brokers,
		Topic:   conf.ConsumeKafka.Topic,
		GroupID: "live-connection-group",
	})

	defer kafkaConsumer.Stop()

	messageProcessor := services.NewKafkaProcessor(hub, kafkaProducer)
	statusStorage := storage.NewOnlineStorage(db)
	statusService := services.NewStatusService(statusStorage, hub)
	statusHandler := handler.NewOnlineStatusServer(statusService)

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
		GET("/health", hub.HealthCheck()).
		GET("/v1.0/status/:users", statusHandler.GetStatus())

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

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
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
