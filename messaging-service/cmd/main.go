package main

import (
	"context"
	"log"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/configuration"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/instrumentation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	config, err := configuration.LoadConfig("/etc/messaging/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	otlpExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(config.Otlp.GrpcAddr),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer otlpExporter.Shutdown(ctx)

	otlpRes, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("backend-services"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(otlpExporter),
		sdktrace.WithResource(otlpRes),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	pgxDB, err := pgxpool.New(ctx, config.DB.ConnString)
	if err != nil {
		log.Fatalf("Connect to pg failed: %s\n", err)
	}
	defer pgxDB.Close()
	db := instrumentation.Tracing(pgxDB)

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis connection establishing failed: %s", err)
	}
	defer rdb.Close()

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Fatalf("Add instrument tracing to redis failed: %s", err)
	}

	fileStConn, err := grpc.NewClient(config.FileStorage.GrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Cannot connect to file storage gRPC: %s", err)
	}
	defer fileStConn.Close()

	confDB := configuration.NewDB(db, rdb)

	confExternal := configuration.NewExternal(fileStConn)

	srv := configuration.NewServices(confDB, confExternal)

	rest := configuration.NewHandlers(srv)

	ginEngine, err := configuration.GinEngine(rest, confDB, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := ginEngine.Run(":5000"); err != nil {
		log.Fatalf("Gin engine running failed: %s", err)
	}
}
