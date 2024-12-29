package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chakchat/chakchat/backend/identity/internal/services" // Actually, I am worried about this dependency
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	prefixPhoneToId = "PhoneToId:"
	prefixMeta      = "SignInMeta:"
)

type SignInMetaConfig struct {
	MetaLifetime time.Duration
}

type SignInMetaStorage struct {
	client *redis.Client
	conf   *SignInMetaConfig
}

func NewSignInMetaStorage(conf *SignInMetaConfig, client *redis.Client) *SignInMetaStorage {
	return &SignInMetaStorage{
		client: client,
		conf:   conf,
	}
}

func (s *SignInMetaStorage) Remove(ctx context.Context, signInKey uuid.UUID) error {
	idKey := prefixMeta + signInKey.String()

	// Here I don't delete PhoneToId relation
	// It may cause bugs if these entities will be used in other scenarios
	// But now PhoneToId just will expire or be overwritten
	res := s.client.Del(ctx, idKey)
	if err := res.Err(); err != nil {
		return err
	}

	return nil
}

func (s *SignInMetaStorage) FindMetaByPhone(ctx context.Context, phone string) (*services.SignInMeta, bool, error) {
	phoneKey := prefixPhoneToId + phone

	idResp := s.client.Get(ctx, phoneKey)
	if err := idResp.Err(); err != nil {
		return nil, false, fmt.Errorf("redis get id by phone failed: %s", err)
	}

	id, err := uuid.Parse(idResp.String())
	if err != nil {
		return nil, false, fmt.Errorf("uuid parsing failed: %s", err)
	}
	return s.FindMeta(ctx, id)
}

func (s *SignInMetaStorage) FindMeta(ctx context.Context, signInKey uuid.UUID) (*services.SignInMeta, bool, error) {
	idKey := prefixMeta + signInKey.String()
	metaResp := s.client.Get(ctx, idKey)
	if err := metaResp.Err(); err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("redis get meta failed: %s", err)
	}

	meta := new(services.SignInMeta)
	if err := json.Unmarshal([]byte(metaResp.String()), meta); err != nil {
		return nil, false, fmt.Errorf("umarshalling meta failed: %s", err)
	}

	return meta, true, nil
}

func (s *SignInMetaStorage) Store(ctx context.Context, meta *services.SignInMeta) error {
	metaJson, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("metadata json marshalling failed: %s", err)
	}

	id := meta.UserId.String()
	metaKey := prefixMeta + id
	phoneKey := prefixPhoneToId + meta.Phone

	status := s.client.Set(ctx, metaKey, metaJson, s.conf.MetaLifetime)
	if err := status.Err(); err != nil {
		return err
	}

	status = s.client.Set(ctx, phoneKey, id, s.conf.MetaLifetime)
	if err := status.Err(); err != nil {
		return err
	}

	return nil
}
