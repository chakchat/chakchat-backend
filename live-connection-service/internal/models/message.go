package models

type KafkaMessage struct {
	Receivers []string `json:"receivers"`
	Type      string   `json:"type"`
	Data      any      `json:"data"`
}

type WSMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
