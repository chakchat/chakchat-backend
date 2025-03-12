package response

import "github.com/chakchat/chakchat-backend/messaging-service/internal/application/dto"

func TextMessage(msg *dto.TextMessageDTO) JSONResponse {
	const (
		ChatIDField    = "chat_id"
		UpdateIDField  = "update_id"
		TypeField      = "type"
		SenderIDField  = "sender_id"
		CreatedAtField = "created_at"

		ContentField  = "content"
		TextField     = "text"
		ReplyToField  = "reply_to"
		EditedIDField = "edited"
	)
	const (
		UpdateTypeTextMessage = "text_message"
	)
	var edited *int64
	if msg.Edited != nil {
		cp := int64(msg.Edited.UpdateID)
		edited = &cp
	}

	return JSONResponse{
		ChatIDField:    msg.ChatID,
		UpdateIDField:  msg.UpdateID,
		TypeField:      UpdateTypeTextMessage,
		SenderIDField:  msg.SenderID,
		CreatedAtField: msg.CreatedAt,

		ContentField: JSONResponse{
			TextField:     msg.Text,
			ReplyToField:  msg.ReplyTo,
			EditedIDField: edited,
		},
	}
}
