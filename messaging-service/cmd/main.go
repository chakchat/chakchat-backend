package main

import (
	"context"
	"log"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/configuration"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := configuration.LoadConfig("/app/config.yml")
	if err != nil {
		log.Fatal(err)
	}

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

	confDB := configuration.NewDB(db, rdb)

	srv := configuration.NewServices(confDB, configuration.NewExternal())

	rest := configuration.NewHandlers(srv)

	ginEngine := configuration.GinEngine(rest, confDB, config)

	if err := ginEngine.Run(":5000"); err != nil {
		log.Fatalf("Gin engine running failed: %s", err)
	}
}
