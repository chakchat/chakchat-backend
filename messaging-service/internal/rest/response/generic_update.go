package response

import (
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/services"
)

func TextMessage(msg *dto.TextMessageDTO) JSONResponse {
	return GenericUpdate(&services.GenericUpdate{
		UpdateID:   msg.UpdateID,
		ChatID:     msg.ChatID,
		SenderID:   msg.SenderID,
		UpdateType: services.UpdateTypeTextMessage,
		CreatedAt:  msg.CreatedAt,
		Info: services.GenericUpdateInfo{
			TextMessage: &services.TextMessageInfo{
				Text:    msg.Text,
				Edited:  convertEdited(msg.Edited),
				ReplyTo: msg.ReplyTo,
			},
		},
	})
}

func TextMessageEdited(edited *dto.TextMessageEditedDTO) JSONResponse {
	return GenericUpdate(convertEdited(edited))
}

func FileMessage(msg *dto.FileMessageDTO) JSONResponse {
	return GenericUpdate(&services.GenericUpdate{
		UpdateID:   msg.UpdateID,
		ChatID:     msg.ChatID,
		SenderID:   msg.SenderID,
		UpdateType: services.UpdateTypeFileMessage,
		CreatedAt:  msg.CreatedAt,
		Info: services.GenericUpdateInfo{
			FileMessage: &services.FileMessageInfo{
				File:    msg.File,
				ReplyTo: msg.ReplyTo,
			},
		},
	})
}

func UpdateDeleted(ud *dto.UpdateDeletedDTO) JSONResponse {
	return GenericUpdate(&services.GenericUpdate{
		UpdateID:   ud.UpdateID,
		ChatID:     ud.ChatID,
		SenderID:   ud.SenderID,
		UpdateType: services.UpdateTypeDeleted,
		CreatedAt:  ud.CreatedAt,
		Info: services.GenericUpdateInfo{
			Deleted: &services.DeletedInfo{
				DeletedID:  ud.DeletedID,
				DeleteMode: ud.DeleteMode,
			},
		},
	})
}

func Reaction(r *dto.ReactionDTO) JSONResponse {
	return GenericUpdate(&services.GenericUpdate{
		UpdateID:   r.UpdateID,
		ChatID:     r.ChatID,
		SenderID:   r.SenderID,
		UpdateType: services.UpdateTypeReaction,
		CreatedAt:  r.CreatedAt,
		Info: services.GenericUpdateInfo{
			Reaction: &services.ReactionInfo{
				Reaction: r.ReactionType,
			},
		},
	})
}

func SecretUpdate(r *dto.SecretUpdateDTO) JSONResponse {
	return GenericUpdate(&services.GenericUpdate{
		UpdateID:   r.UpdateID,
		ChatID:     r.ChatID,
		SenderID:   r.SenderID,
		UpdateType: services.UpdateTypeSecret,
		CreatedAt:  r.CreatedAt,
		Info: services.GenericUpdateInfo{
			Secret: &services.SecretUpdateInfo{
				Payload:              r.Payload,
				InitializationVector: r.InitializationVector,
				KeyHash:              r.KeyHash,
			},
		},
	})
}

func GenericUpdate(update *services.GenericUpdate) JSONResponse {
	const (
		ChatIDField    = "chat_id"
		UpdateIDField  = "update_id"
		TypeField      = "type"
		SenderIDField  = "sender_id"
		CreatedAtField = "created_at"

		ContentField = "content"

		TextField    = "text"
		ReplyToField = "reply_to"
		EditedField  = "edited"

		NewTextField   = "new_text"
		MessageIDField = "message_id"

		ReactionsField = "reactions"

		FileField = "file"

		FileIDField       = "file_id"
		FileNameField     = "file_name"
		FileMimetypeField = "mime_type"
		FileSizeField     = "file_size"
		FileURLField      = "file_url"

		DeletedIDField   = "deleted_id"
		DeletedModeField = "deleted_mode"

		ReactionField = "reaction"

		SecretPayload              = "payload"
		SecretInitializationVector = "initialization_vector"
		SecretKeyHash              = "key_hash"
	)

	resp := JSONResponse{
		ChatIDField:    update.ChatID,
		UpdateIDField:  update.UpdateID,
		TypeField:      update.UpdateType,
		SenderIDField:  update.SenderID,
		CreatedAtField: update.CreatedAt,
	}

	switch update.UpdateType {
	case services.UpdateTypeTextMessage:
		resp[ContentField] = JSONResponse{
			TextField:    update.Info.TextMessage.Text,
			ReplyToField: update.Info.TextMessage.ReplyTo,
			ReactionsField: GenericUpdates(update.Info.TextMessage.Reactions),
		}
		if update.Info.TextMessage.Edited != nil {
			resp[EditedField] = GenericUpdate(update.Info.TextMessage.Edited)
		}
	case services.UpdateTypeTextMessageEdited:
		resp[ContentField] = JSONResponse{
			NewTextField:   update.Info.TextMessageEdited.NewText,
			MessageIDField: update.Info.TextMessageEdited.MessageID,
		}
	case services.UpdateTypeFileMessage:
		resp[ContentField] = JSONResponse{
			FileField: JSONResponse{
				FileIDField:       update.Info.FileMessage.File.FileId,
				FileNameField:     update.Info.FileMessage.File.FileName,
				FileMimetypeField: update.Info.FileMessage.File.MimeType,
				FileSizeField:     update.Info.FileMessage.File.FileSize,
				FileURLField:      update.Info.FileMessage.File.FileURL,
				CreatedAtField:    update.Info.FileMessage.File.CreatedAt,
			},
			ReplyToField: update.Info.FileMessage.ReplyTo,
			ReactionsField: GenericUpdates(update.Info.FileMessage.Reactions),
		}
	case services.UpdateTypeDeleted:
		resp[ContentField] = JSONResponse{
			DeletedIDField:   update.Info.Deleted.DeletedID,
			DeletedModeField: update.Info.Deleted.DeleteMode,
		}
	case services.UpdateTypeReaction:
		resp[ContentField] = JSONResponse{
			ReactionField: update.Info.Reaction.Reaction,
		}
	case services.UpdateTypeSecret:
		resp[ContentField] = JSONResponse{
			SecretKeyHash:              update.Info.Secret.KeyHash,
			SecretInitializationVector: update.Info.Secret.InitializationVector,
			SecretPayload:              update.Info.Secret.Payload,
		}
	}

	return resp
}

func GenericUpdates(updates []services.GenericUpdate) []JSONResponse {
	res := make([]JSONResponse, len(updates))
	for i, up := range updates {
		res[i] = GenericUpdate(&up)
	}
	return res
}

func convertEdited(edited *dto.TextMessageEditedDTO) *services.GenericUpdate {
	if edited == nil {
		return nil
	}

	return &services.GenericUpdate{
		UpdateID:   edited.UpdateID,
		ChatID:     edited.ChatID,
		SenderID:   edited.SenderID,
		UpdateType: services.UpdateTypeTextMessageEdited,
		CreatedAt:  edited.CreatedAt,
		Info: services.GenericUpdateInfo{
			TextMessageEdited: &services.TextMessageEditedInfo{
				MessageID: edited.MessageID,
				NewText:   edited.NewText,
			},
		},
	}
}
