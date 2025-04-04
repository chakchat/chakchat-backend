package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/chat"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	PersonalChat       repository.PersonalChatRepository
	GroupChat          repository.GroupChatRepository
	SecretPersonalChat repository.SecretPersonalChatRepository
	SecretGroupChat    repository.SecretGroupChatRepository
	Chatter            repository.ChatterRepository
	GenericChat        chat.GenericChatRepository

	SQLer storage.SQLer

	Redis *redis.Client
}

func NewDB(db storage.SQLer, redis *redis.Client) *DB {
	return &DB{
		PersonalChat:       chat.NewPersonalChatRepository(),
		GroupChat:          chat.NewGroupChatRepository(),
		SecretPersonalChat: chat.NewSecretPersonalChatRepository(),
		SecretGroupChat:    chat.NewSecretGroupChatRepository(),
		Chatter:            chat.NewChatterRepository(),
		GenericChat:        *chat.NewGenericChatRepository(),
		SQLer:              db,
		Redis:              redis,
	}
}
