package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SmsServerStubSender struct {
	addr string
}

func NewSmsServerStubSender(addr string) *SmsServerStubSender {
	return &SmsServerStubSender{
		addr: addr,
	}
}

func (s *SmsServerStubSender) SendSms(ctx context.Context, phone string, message string) error {
	type Req struct {
		Phone   string `json:"phone"`
		Message string `json:"message"`
	}
	req, _ := json.Marshal(Req{
		Phone:   phone,
		Message: message,
	})
	resp, err := http.Post(s.addr, "application/json", bytes.NewReader(req))
	if err != nil {
		return fmt.Errorf("sending sms to stub server failed: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("sending sms request failed with status code: %s", resp.Status)
	}
	return nil
}

type SmsSenderFake struct{}

func (s *SmsSenderFake) SendSms(ctx context.Context, phone string, message string) error {
	fmt.Printf("Sent sms to %s. Sms message: \"%s\"", phone, message)
	return nil
}
