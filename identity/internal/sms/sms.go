package sms

import (
	"context"
	"fmt"
)

type SmsSenderFake struct{}

func (s *SmsSenderFake) SendSms(ctx context.Context, phone string, message string) error {
	fmt.Printf("Sent sms to %s. Sms message: \"%s\"", phone, message)
	return nil
}
