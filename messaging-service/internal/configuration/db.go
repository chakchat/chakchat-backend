package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
)

type DB struct {
	PersonalChatRepo repository.PersonalChatRepository
}

func NewDB(db *pgx.Conn) *DB {
	return &DB{
		PersonalChatRepo: postgres.NewPersonalChatRepository(db),
	}
}
