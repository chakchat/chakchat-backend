package idempotency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	prefixIdempotencyMeta = "IdempotencyData:Meta"
	prefixIdempotencyData = "IdempotencyData:Data"
)

type IdempotencyConfig struct {
	DataExp time.Duration
}

type idempotencyStorage struct {
	conf   *IdempotencyConfig
	client *redis.Client
}

func NewStorage(client *redis.Client, conf *IdempotencyConfig) IdempotencyStorage {
	return &idempotencyStorage{
		conf:   conf,
		client: client,
	}
}

type cachedRespMeta struct {
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	BodyKey    string      `json:"body_key"`
}

func (s *idempotencyStorage) Get(ctx context.Context, key string) (*CapturedResponse, bool, error) {
	metaKey := prefixIdempotencyMeta + key

	metaRes := s.client.Get(ctx, metaKey)
	if err := metaRes.Err(); err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("idempotency storage's get cached response failed: %s", err)
	}

	var meta cachedRespMeta
	if err := json.Unmarshal([]byte(metaRes.String()), &meta); err != nil {
		return nil, false, fmt.Errorf("cached response metadata unmarshalling failed: %s", err)
	}

	body, err := s.getBody(ctx, meta.BodyKey)
	if err != nil {
		return nil, false, fmt.Errorf("redis getting response failed: %s", err)
	}

	resp := &CapturedResponse{
		StatusCode: meta.StatusCode,
		Headers:    meta.Headers,
		Body:       body,
	}
	return resp, false, nil
}

func (s *idempotencyStorage) getBody(ctx context.Context, bodyKey string) ([]byte, error) {
	bodyRes := s.client.Get(ctx, bodyKey)
	if err := bodyRes.Err(); err != nil {
		return nil, fmt.Errorf("getting body failed: %s", err)
	}

	body, err := bodyRes.Bytes()
	if err != nil {
		return nil, fmt.Errorf("getting body failed: %s", err)
	}
	return body, nil
}

func (s *idempotencyStorage) Store(ctx context.Context, key string, resp *CapturedResponse) error {
	meta := cachedRespMeta{
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		BodyKey:    prefixIdempotencyData + uuid.NewString(),
	}

	metaKey := prefixIdempotencyMeta + meta.BodyKey
	if err := s.client.Set(ctx, metaKey, meta, s.conf.DataExp).Err(); err != nil {
		return fmt.Errorf("redis response caching failed: %s", err)
	}

	if err := s.client.Set(ctx, meta.BodyKey, resp.Body, s.conf.DataExp).Err(); err != nil {
		return fmt.Errorf("redis response caching failed: %s", err)
	}

	return nil
}
