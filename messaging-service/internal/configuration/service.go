package configuration

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/chat"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services/update"
)

type Services struct {
	PersonalChat       *chat.PersonalChatService
	GroupChat          *chat.GroupChatService
	GroupPhoto         *chat.GroupPhotoService
	SecretPersonalChat *chat.SecretPersonalChatService
	SecretGroup        *chat.SecretGroupChatService
	SecretGroupPhoto   *chat.SecretGroupPhotoService
	GenericChat        *chat.GenericChatService

	PersonalUpdate       *update.PersonalUpdateService
	PersonalFile         *update.PersonalFileService
	GroupUpdate          *update.GroupUpdateService
	GroupFile            *update.GroupFileService
	SecretPersonalUpdate *update.SecretPersonalUpdateService
	SecretGroupUpdate    *update.SecretGroupUpdateService
	GenericUpdate        *update.GenericUpdateService
}

func NewServices(db *DB, external *External) *Services {
	return &Services{
		PersonalChat: chat.NewPersonalChatService(
			db.SQLer, db.PersonalChat, external.Publisher,
		),
		GroupChat: chat.NewGroupChatService(
			db.SQLer, db.GroupChat, external.Publisher,
		),
		GroupPhoto: chat.NewGroupPhotoService(
			db.SQLer, db.GroupChat, external.FileStorage, external.Publisher,
		),
		SecretPersonalChat: chat.NewSecretPersonalChatService(
			db.SQLer, db.SecretPersonalChat, external.Publisher,
		),
		SecretGroup: chat.NewSecretGroupChatService(
			db.SQLer, db.SecretGroupChat, external.Publisher,
		),
		SecretGroupPhoto: chat.NewSecretGroupPhotoService(
			db.SQLer, db.SecretGroupChat, external.FileStorage, external.Publisher,
		),
		GenericChat: chat.NewGenericChatService(
			db.SQLer, db.GenericChat,
		),
		PersonalUpdate: update.NewPersonalUpdateService(
			db.SQLer, db.PersonalChat, db.Update, db.Chatter, external.Publisher,
		),
		PersonalFile: update.NewPersonalFileService(
			db.SQLer, db.PersonalChat, db.Update, external.FileStorage, external.Publisher,
		),
		GroupUpdate: update.NewGroupUpdateService(
			db.SQLer, db.GroupChat, db.Update, db.Chatter, external.Publisher,
		),
		GroupFile: update.NewGroupFileService(
			db.SQLer, db.GroupChat, db.Update, external.FileStorage, external.Publisher,
		),
		SecretPersonalUpdate: update.NewSecretPersonalUpdateService(
			db.SQLer, db.SecretPersonalChat, db.SecretUpdate, external.Publisher,
		),
		SecretGroupUpdate: update.NewSecretGroupUpdateService(
			db.SQLer, db.SecretGroupChat, db.SecretUpdate, external.Publisher,
		),
		GenericUpdate: update.NewGenericUpdateService(
			db.SQLer, db.Chatter, db.GenericUpdate,
		),
	}
}
