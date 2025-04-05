package models

type KafkaMessageUpdate struct {
	Receivers []string       `json:"receivers"`
	Type      string         `json:"type"`
	Data      UpdateResponse `json:"data"`
}

type WSMessageUpdate struct {
	Type string         `json:"type"`
	Data UpdateResponse `json:"data"`
}

type UpdateResponse struct {
	ChatID    string         `json:"chat_id"`
	UpdateID  string         `json:"update_id"`
	Type      string         `json:"type"`
	SenderID  string         `json:"sender_id"`
	CreatedAt string         `json:"created_at"`
	Content   map[string]any `json:"content"`
}

type TextMessageContent struct {
	Text    string    `json:"text"`
	ReplyTo *string   `json:"reply_to,omitempty"`
	Edited  *EditInfo `json:"edited,omitempty"`
}

type FileMessageContent struct {
	File    FileInfo `json:"file"`
	ReplyTo *string  `json:"reply_to,omitempty"`
}

type EditInfo struct {
	NewText   string `json:"new_text"`
	MessageID string `json:"message_id"`
}

type FileInfo struct {
	FileID    string `json:"file_id"`
	FileName  string `json:"file_name"`
	MimeType  string `json:"mime_type"`
	FileSize  int64  `json:"file_size"`
	FileURL   string `json:"file_url"`
	CreatedAt string `json:"created_at"`
}
