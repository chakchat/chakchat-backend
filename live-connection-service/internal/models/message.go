package models

import "github.com/google/uuid"

type KafkaMessage struct {
	Receivers []uuid.UUID `json:"receivers"`
	Type      string      `json:"type"`
	Data      any         `json:"data"`
}

type WSMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
