package services

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/domain/models"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/infrastructure/messages"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/infrastructure/ws"
)

type KafkaProcessor struct {
	hub *ws.Hub
	dlq *messages.Producer
}

func (p *KafkaProcessor) MessageHandler(ctx context.Context, msg kafka.Message) {
	msgType, err := p.detectMessageType(msg.Value)
	if err != nil {
		p.dlq.Send(ctx, msg.Value)
	}

	if msgType == "update" {
		p.processUpdateMessage(ctx, msg.Value)
	} else {
		p.processChatMessage(ctx, msg.Value)
	}
}

func (p *KafkaProcessor) detectMessageType(data []byte) (string, error) {
	var msg struct {
		Receivers []string `json:"receivers"`
		Type      string   `json:"type"`
		Data      any      `json:"data"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return "", err
	}

	return msg.Type, nil
}

func (p *KafkaProcessor) processUpdateMessage(ctx context.Context, data []byte) {
	var update models.KafkaMessageUpdate
	if err := json.Unmarshal(data, &update); err != nil {
		return
	}

	response := models.WSMessageUpdate{
		Type: update.Type,
		Data: update.Data,
	}
	for _, userId := range update.Receivers {
		if !p.hub.Send(userId, response) {
			p.dlq.Send(ctx, data)
		}
	}

}

func (p *KafkaProcessor) processChatMessage(ctx context.Context, data []byte) {
	var update models.KafkaMessageChat
	if err := json.Unmarshal(data, &update); err != nil {
		return
	}

	response := models.WSMessageChat{
		Type: update.Type,
		Data: update.Data,
	}
	for _, userId := range update.Receivers {
		if !p.hub.Send(userId, response) {
			p.dlq.Send(ctx, data)
		}
	}

}
