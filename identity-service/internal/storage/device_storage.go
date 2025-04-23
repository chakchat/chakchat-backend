package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chakchat/chakchat-backend/identity-service/internal/services"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const DeviceKeyPrefix = "Device:User:"

type DeviceStorageConfig struct {
	DeviceInfoLifetime time.Duration
}

type DeviceStorage struct {
	client *redis.Client
	config *DeviceStorageConfig
}

func NewDeviceStorage(client *redis.Client, config *DeviceStorageConfig) *DeviceStorage {
	return &DeviceStorage{
		client: client,
		config: config,
	}
}

func (s *DeviceStorage) Store(ctx context.Context, userID uuid.UUID, info *services.DeviceInfo) error {
	key := DeviceKeyPrefix + userID.String()

	enc, err := json.Marshal(info)
	if err != nil {
		return err
	}

	status := s.client.Set(ctx, key, enc, s.config.DeviceInfoLifetime)
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

func (s *DeviceStorage) Refresh(ctx context.Context, userID uuid.UUID) error {
	key := DeviceKeyPrefix + userID.String()
	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	if exists == 1 {
		err := s.client.Expire(ctx, key, s.config.DeviceInfoLifetime).Err()
		return err
	}
	return nil
}

func (s *DeviceStorage) Remove(ctx context.Context, userID uuid.UUID) error {
	key := DeviceKeyPrefix + userID.String()
	res := s.client.Del(ctx, key)
	if err := res.Err(); err != nil {
		return err
	}
	return nil
}
