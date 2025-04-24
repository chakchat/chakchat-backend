package notifier

import "github.com/sideshow/apns2"

type APNSClient interface {
	SendNotification(deviceToken, title string) error
}

type APNsClient struct {
	client *apns2.Client
}

func NewAPNsClient(certPath, keyID, teamID string) *APNSClient {
	return nil
}

func (a *APNsClient) SendNotification(deviceToken *string, title *string) error {
	return nil
}
