package notifier

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

type APNSClient interface {
	SendNotification(receiver, title string) error
}

type APNsClient struct {
	client *apns2.Client
	topic  string
}

func NewAPNsClient(certPath, keyID, teamID, topic string) (*APNsClient, error) {
	authKey, err := token.AuthKeyFromFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load APNs certificate: %w", err)
	}

	token := &token.Token{
		AuthKey: authKey,
		KeyID:   keyID,
		TeamID:  teamID,
	}
	client := apns2.NewTokenClient(token).Development()

	return &APNsClient{client: client, topic: topic}, nil
}

func (a *APNsClient) SendNotification(deviceToken, receiver, title string) error {
	payload := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": map[string]string{
				"title": receiver,
				"body":  title,
			},
			"sound": "default",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	notification := &apns2.Notification{
		DeviceToken: deviceToken,
		Topic:       a.topic,
		Payload:     payloadBytes,
	}

	response, err := a.client.Push(notification)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	if !response.Sent() {
		return fmt.Errorf("failed to send notification: %v", response.Reason)
	}

	log.Printf("Notification sent to %s: %s", deviceToken, title)
	return nil
}
