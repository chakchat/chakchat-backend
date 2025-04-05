package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage/repository"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/chat"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/update"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	PersonalChat       repository.PersonalChatRepository
	GroupChat          repository.GroupChatRepository
	SecretPersonalChat repository.SecretPersonalChatRepository
	SecretGroupChat    repository.SecretGroupChatRepository
	Chatter            repository.ChatterRepository
	GenericChat        repository.GenericChatRepository

	UpdateRepository       repository.UpdateRepository
	SecretUpdateRepository repository.SecretUpdateRepository

	SQLer storage.SQLer

	Redis *redis.Client
}

func NewDB(db storage.SQLer, redis *redis.Client) *DB {
	return &DB{
		PersonalChat:           chat.NewPersonalChatRepository(),
		GroupChat:              chat.NewGroupChatRepository(),
		SecretPersonalChat:     chat.NewSecretPersonalChatRepository(),
		SecretGroupChat:        chat.NewSecretGroupChatRepository(),
		Chatter:                chat.NewChatterRepository(),
		GenericChat:            chat.NewGenericChatRepository(),
		UpdateRepository:       update.NewUpdateRepository(),
		SecretUpdateRepository: update.NewSecretUpdateRepository(),
		SQLer:                  db,
		Redis:                  redis,
	}
}
