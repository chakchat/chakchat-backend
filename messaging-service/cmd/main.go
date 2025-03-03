package main

import (
	"context"
	"log"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/configuration"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	config, err := configuration.LoadConfig("/app/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	otlpExporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(config.Otlp.GrpcAddr),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer otlpExporter.Shutdown(context.Background())

	otlpRes, err := resource.New(
		context.Background(),
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

	db, err := pgx.Connect(context.Background(), config.DB.ConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis connection establishing failed: %s", err)
	}
	defer rdb.Close()

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Fatalf("Add instrument tracing to redis failed: %s", err)
	}

	confDB := configuration.NewDB(db, rdb)

	srv := configuration.NewServices(confDB, configuration.NewExternal())

	rest := configuration.NewHandlers(srv)

	ginEngine := configuration.GinEngine(rest, confDB, config)

	if err := ginEngine.Run(":5000"); err != nil {
		log.Fatalf("Gin engine running failed: %s", err)
	}
}
