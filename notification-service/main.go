package main

import (
	"context"
	"log"

	"github.com/chakchat/chakchat-backend/notification-service/internal/identity"
	"github.com/chakchat/chakchat-backend/notification-service/internal/user"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Identity struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"identity"`

	User struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"user"`

	Otlp struct {
		GrpcAddr string `mapstructure:"grpc_addr"`
	} `mapstructure:"otlp"`

	APNs struct {
		CertPath string `mapstructure:"cert_path"`
		KeyPath  string `mapstructure:"key_path"`
		TeamId   string `mapstructure:"team_id"`
	} `mapstructure:"apns"`
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

var conf *Config = loadConfig("/app/config.yml")

func main() {
	// identityClient, closeIdentity := createIdentityClient()
	// userClient, closeUser := createUserClient()

	// defer closeIdentity()
	// defer closeUser()

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

	// grpcService := grpc_service.NewGrpcClients(userClient, identityClient)
	// parser := notifier.NewParser(grpcService)
	// apnsClient := notifier.NewAPNsClient(conf.APNs.CertPath, conf.APNs.KeyPath, conf.APNs.TeamId)
	// parser.ParseNotification(context.Background(), []byte("something"))

	// userID, err := uuid.Parse("ytjdjdruktdkfc")
	// if err != nil {
	// 	log.Printf("Invalid UUID format: %v", err)
	// 	return
	// }

}

func createIdentityClient() (identity.IdentityServiceClient, func() error) {
	addr := conf.Identity.GrpcAddr
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return identity.NewIdentityServiceClient(conn), conn.Close
}

func createUserClient() (user.UserServiceClient, func() error) {
	addr := conf.Identity.GrpcAddr
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return user.NewUserServiceClient(conn), conn.Close
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
