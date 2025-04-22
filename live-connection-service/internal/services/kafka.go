package services

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"

	"github.com/chakchat/chakchat-backend/live-connection-service/internal/models"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/mq"
	"github.com/chakchat/chakchat-backend/live-connection-service/internal/ws"
)

type KafkaProcessor struct {
	hub    *ws.Hub
	notifq *mq.Producer //queue to send message in notification service
}

func NewKafkaProcessor(hub *ws.Hub, dlq *mq.Producer) *KafkaProcessor {
	return &KafkaProcessor{
		hub:    hub,
		notifq: dlq,
	}
}

func (p *KafkaProcessor) MessageHandler(ctx context.Context, msg kafka.Message) error {
	var message models.KafkaMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return err
	}

	response := models.WSMessage{
		Type: message.Type,
		Data: message.Data,
	}

	var notificReceivers []string
	for _, userId := range message.Receivers {
		if !p.hub.Send(userId, response) {
			notificReceivers = append(notificReceivers, userId)
		}
	}
	if len(notificReceivers) != 0 {
		notificMessage := models.KafkaMessage{
			Receivers: notificReceivers,
			Type:      message.Type,
			Data:      message.Data,
		}
		p.notifq.Send(ctx, notificMessage)
	}
	return nil
}
