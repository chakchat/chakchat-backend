package response

import (
	"fmt"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
)

type JSONResponse = map[string]any

func PersonalChat(chat *dto.PersonalChatDTO) JSONResponse {
	generic := services.NewPersonalGenericChat(chat.ID, chat.CreatedAt, chat.Members[:], services.PersonalInfo{
		BlockedBy: chat.BlockedBy,
	})
	return GenericChat(&generic)
}

func GroupChat(chat *dto.GroupChatDTO) JSONResponse {
	generic := services.NewGroupGenericChat(chat.ID, chat.CreatedAt, chat.Members, services.GroupInfo{
		Admin:            chat.Admin,
		GroupName:        chat.Name,
		GroupDescription: chat.Description,
		GroupPhoto:       chat.GroupPhoto,
	})
	return GenericChat(&generic)
}

func SecretPersonalChat(chat *dto.SecretPersonalChatDTO) JSONResponse {
	generic := services.NewSecretPersonalGenericChat(chat.ID, chat.CreatedAt, chat.Members[:], services.SecretPersonalInfo{
		Expiration: chat.Expiration,
	})
	return GenericChat(&generic)
}

func SecretGroupChat(chat *dto.SecretGroupChatDTO) JSONResponse {
	generic := services.NewSecretGroupGenericChat(chat.ID, chat.CreatedAt, chat.Members, services.SecretGroupInfo{
		Admin:            chat.Admin,
		GroupName:        chat.Name,
		GroupDescription: chat.Description,
		GroupPhoto:       chat.GroupPhotoURL,
		Expiration:       chat.Expiration,
	})
	return GenericChat(&generic)
}

func GenericChat(chat *services.GenericChat) JSONResponse {
	const (
		ChatIDField    = "chat_id"
		TypeField      = "type"
		MembersField   = "members"
		CreatedAtField = "created_at"
		InfoField      = "info"

		LastUpdateIDField  = "last_update_id"
		UpdatePreviewField = "update_preview"

		BlockedByField = "blocked_by"

		AdminField       = "admin_id"
		NameField        = "name"
		DescriptionField = "description"
		GroupPhotoField  = "group_photo"

		ExpirationField = "expiration"
	)

	resp := JSONResponse{
		ChatIDField:    chat.ChatID,
		TypeField:      chat.ChatType,
		MembersField:   chat.Members,
		CreatedAtField: chat.CreatedAt,
	}
	if chat.LastUpdateID != nil {
		resp[LastUpdateIDField] = *chat.LastUpdateID
	}
	if chat.UpdatePreview != nil {
		updates := make([]JSONResponse, len(chat.UpdatePreview))
		for i := range chat.UpdatePreview {
			updates[i] = GenericUpdate(&chat.UpdatePreview[i])
		}
		resp[UpdatePreviewField] = updates
	}

	switch chat.ChatType {
	case services.ChatTypePersonal:
		resp[InfoField] = JSONResponse{
			BlockedByField: chat.Info.Personal().BlockedBy,
		}
	case services.ChatTypeGroup:
		resp[InfoField] = JSONResponse{
			AdminField:       chat.Info.Group().Admin,
			NameField:        chat.Info.Group().GroupName,
			DescriptionField: chat.Info.Group().GroupDescription,
			GroupPhotoField:  chat.Info.Group().GroupPhoto,
		}
	case services.ChatTypeSecretPersonal:
		resp[InfoField] = JSONResponse{
			ExpirationField: chat.Info.SecretPersonal().Expiration,
		}
	case services.ChatTypeSecretGroup:
		resp[InfoField] = JSONResponse{
			AdminField:       chat.Info.SecretGroup().Admin,
			NameField:        chat.Info.SecretGroup().GroupName,
			DescriptionField: chat.Info.SecretGroup().GroupDescription,
			GroupPhotoField:  chat.Info.SecretGroup().GroupPhoto,
			ExpirationField:  chat.Info.SecretPersonal().Expiration,
		}
	default:
		panic(fmt.Errorf("uknown chat type: %s", chat.ChatType))
	}

	return resp
}
