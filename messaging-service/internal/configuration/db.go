package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/chat"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	PersonalChat repository.PersonalChatRepository
	Redis        *redis.Client
}

func NewDB(db *pgx.Conn, redis *redis.Client) *DB {
	return &DB{
		PersonalChat: chat.NewPersonalChatRepository(db),
		Redis:        redis,
	}
}
