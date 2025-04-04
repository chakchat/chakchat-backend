package models

type KafkaMessageChat struct {
	Receivers []string     `json:"receivers"`
	Type      string       `json:"type"`
	Data      ChatResponse `json:"data"`
}

type WSMessageChat struct {
	Type string       `json:"type"`
	Data ChatResponse `json:"data"`
}

type ChatResponse struct {
	ChatID    string         `json:"chat_id"`
	Type      string         `json:"type"`
	Members   []string       `json:"members"`
	CreatedAt string         `json:"created_at"`
	Info      map[string]any `json:"info"`
}

type PersonalChatInfo struct {
	BlockedBy []string `json:"blocked_by"`
}

type GroupChatInfo struct {
	Admin       string `json:"admin_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupPhoto  string `json:"group_photo"`
}

type SecretChatInfo struct {
	Expiration string `json:"expiration"`
}
