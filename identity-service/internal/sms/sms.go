package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SmsSender interface {
	SendSms(ctx context.Context, phone string, message string) error
}

type Client struct {
	email   string
	apiKey  string
	from    string
	client  *http.Client
	baseUrl string
}

type SmsAeroRequest struct {
	method   string
	endpoint string
	params   url.Values
	body     any
}

type SmsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID     int    `json:"id"`
		From   string `json:"from"`
		Number string `json:"number"`
		Text   string `json:"text"`
		Status string `json:"status"`
		ExtID  string `json:"ext_id"`
	} `json:"data"`
	Message string `json:"message"`
}

func NewSmsSender(email string, apiKey string, from string) *Client {
	return &Client{
		email:   email,
		apiKey:  apiKey,
		from:    from,
		client:  &http.Client{},
		baseUrl: "https://gate.smsaero.ru/v2",
	}
}

func (s *Client) executeRequest(ctx context.Context, req SmsAeroRequest) ([]byte, error) {
	reqUrl := fmt.Sprintf("%s/%s?%s", s.baseUrl, req.endpoint, req.params.Encode())

	var reqBody io.Reader
	if req.body != nil {
		jsonBody, err := json.Marshal(req.body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	newReq, err := http.NewRequestWithContext(ctx, req.method, reqUrl, reqBody)
	if err != nil {
		return nil, err
	}

	newReq.SetBasicAuth(s.email, s.apiKey)
	newReq.Header.Set("Accept", "application/json")
	if req.body != nil {
		newReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(newReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	return respBody, nil
}

func (s *Client) SendSms(ctx context.Context, phone string, message string) error {
	params := url.Values{}
	params.Add("number", phone)
	params.Add("text", message)
	params.Add("sign", s.from)
	params.Add("channel", "DIGITAL")

	resp, err := s.executeRequest(ctx, SmsAeroRequest{
		method:   http.MethodPost,
		endpoint: "sms/send",
		params:   params,
		body:     nil,
	})

	if err != nil {
		return err
	}

	var smsResp SmsResponse
	if err := json.Unmarshal(resp, &smsResp); err != nil {
		return err
	}
	return nil
}

type SmsSenderFake struct{}

func (s *SmsSenderFake) SendSms(ctx context.Context, phone string, message string) error {
	fmt.Printf("Sent sms to %s. Sms message: \"%s\"", phone, message)
	return nil
}
