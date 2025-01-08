package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chakchat/chakchat/backend/identity/internal/services" // Actually, I am worried about this dependency
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	prefixPhoneToSignInKey = "PhoneToKey:"
	prefixSignInMeta       = "SignInMeta:"
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
	idKey := prefixSignInMeta + signInKey.String()

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
	phoneKey := prefixPhoneToSignInKey + phone

	keyResp := s.client.Get(ctx, phoneKey)
	if err := keyResp.Err(); err != nil {
		if err == redis.Nil {
			log.Printf("phone-to-key not found in redis: %s", phoneKey)
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("redis get key by phone failed: %s", err)
	}

	key, err := uuid.Parse(keyResp.Val())
	if err != nil {
		log.Printf("uuid parsing failed. uuid was: %s", keyResp.Val())
		return nil, false, fmt.Errorf("uuid parsing failed: %s", err)
	}
	return s.FindMeta(ctx, key)
}

func (s *SignInMetaStorage) FindMeta(ctx context.Context, signInKey uuid.UUID) (*services.SignInMeta, bool, error) {
	idKey := prefixSignInMeta + signInKey.String()
	metaResp := s.client.Get(ctx, idKey)
	if err := metaResp.Err(); err != nil {
		if err == redis.Nil {
			log.Printf("meta not found in redis: %s", idKey)
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("redis get meta failed: %s", err)
	}

	meta := new(services.SignInMeta)
	if err := json.Unmarshal([]byte(metaResp.Val()), meta); err != nil {
		return nil, false, fmt.Errorf("umarshalling meta failed: %s", err)
	}

	return meta, true, nil
}

func (s *SignInMetaStorage) Store(ctx context.Context, meta *services.SignInMeta) error {
	metaJson, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("metadata json marshalling failed: %s", err)
	}

	key := meta.SignInKey.String()
	metaKey := prefixSignInMeta + key
	phoneKey := prefixPhoneToSignInKey + meta.Phone

	status := s.client.Set(ctx, metaKey, metaJson, s.conf.MetaLifetime)
	if err := status.Err(); err != nil {
		return err
	}
	log.Printf("meta stored in redis: key=%s, meta=%v", metaKey, meta)

	status = s.client.Set(ctx, phoneKey, key, s.conf.MetaLifetime)
	if err := status.Err(); err != nil {
		return err
	}
	log.Printf("phone-to-key stored in redis: phoneKey=%s, key=%v", phoneKey, key)

	return nil
}
