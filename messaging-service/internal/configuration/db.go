package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/chat"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	PersonalChat       repository.PersonalChatRepository
	GroupChat          repository.GroupChatRepository
	SecretPersonalChat repository.SecretPersonalChatRepository
	SecretGroupChat    repository.SecretGroupChatRepository
	Chatter            repository.ChatterRepository

	TxProvider storage.TxProvider

	Redis *redis.Client
}

func NewDB(db *pgx.Conn, redis *redis.Client) *DB {
	return &DB{
		PersonalChat:       chat.NewPersonalChatRepository(db),
		GroupChat:          chat.NewGroupChatRepository(db),
		SecretPersonalChat: chat.NewSecretPersonalChatRepository(db),
		SecretGroupChat:    chat.NewSecretGroupChatRepository(db),
		Chatter:            chat.NewChatterRepository(db),
		TxProvider:         postgres.TxProviderStub{},
		Redis:              redis,
	}
}
