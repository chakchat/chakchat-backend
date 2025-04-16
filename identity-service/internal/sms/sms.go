package sms

import (
	"context"
	"fmt"
	"strconv"

	smsaero_golang "github.com/smsaero/smsaero_golang/smsaero"
)

type Client struct {
	email  string
	apiKey string
}
type SmsSender interface {
	SendSms(context.Context, string, string) (*smsaero_golang.SendSms, error)
}

func NewSmsSender(email string, apiKey string) *Client {
	return &Client{
		email:  email,
		apiKey: apiKey,
	}
}

func (c *Client) SendSms(ctx context.Context, phone string, sms string) (*smsaero_golang.SendSms, error) {
	client := smsaero_golang.NewSmsAeroClient(c.email, c.apiKey, smsaero_golang.WithContext(ctx))
	phoneInt, err := strconv.Atoi(phone)
	if err != nil {
		return nil, fmt.Errorf("error converting phone number to integer: %v", err)
	}
	message, err := client.SendSms(phoneInt, sms)
	if err != nil {
		return nil, err
	}
	return &message, nil
}
