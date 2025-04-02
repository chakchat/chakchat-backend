package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/chat"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/tx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	PersonalChat       repository.PersonalChatRepository
	GroupChat          repository.GroupChatRepository
	SecretPersonalChat repository.SecretPersonalChatRepository
	SecretGroupChat    repository.SecretGroupChatRepository
	Chatter            repository.ChatterRepository
	GenericChat        chat.GenericChatRepository

	TxProvider storage.TxProvider

	Redis *redis.Client
}

func NewDB(db *pgxpool.Pool, redis *redis.Client) *DB {
	return &DB{
		PersonalChat:       chat.NewPersonalChatRepository(),
		GroupChat:          chat.NewGroupChatRepository(),
		SecretPersonalChat: chat.NewSecretPersonalChatRepository(),
		SecretGroupChat:    chat.NewSecretGroupChatRepository(),
		Chatter:            chat.NewChatterRepository(),
		GenericChat:        *chat.NewGenericChatRepository(),
		TxProvider:         tx.NewTxProvider(db),
		Redis:              redis,
	}
}
