package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chakchat/chakchat/backend/identity/internal/services"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	prefixPhoneToSignUpKey = "PhoneToSignUpKey:"
	prefixSignUpMeta       = "SignUpMeta:"
)

type SignUpMetaConfig struct {
	MetaLifetime time.Duration
}

type SignUpMetaStorage struct {
	client *redis.Client
	conf   *SignUpMetaConfig
}

func NewSignUpMetaStorage(conf *SignUpMetaConfig, client *redis.Client) *SignUpMetaStorage {
	return &SignUpMetaStorage{
		client: client,
		conf:   conf,
	}
}

func (s *SignUpMetaStorage) Remove(ctx context.Context, signUpKey uuid.UUID) error {
	idKey := prefixSignUpMeta + signUpKey.String()

	// Here I don't delete PhoneToId relation
	// It may cause bugs if these entities will be used in other scenarios
	// But now PhoneToId just will expire or be overwritten
	res := s.client.Del(ctx, idKey)
	if err := res.Err(); err != nil {
		return err
	}

	return nil
}

func (s *SignUpMetaStorage) FindMetaByPhone(ctx context.Context, phone string) (*services.SignUpMeta, bool, error) {
	phoneKey := prefixPhoneToSignUpKey + phone

	keyResp := s.client.Get(ctx, phoneKey)
	if err := keyResp.Err(); err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("redis get sign up key by phone failed: %s", err)
	}

	key, err := uuid.Parse(keyResp.Val())
	if err != nil {
		return nil, false, fmt.Errorf("uuid parsing failed: %s", err)
	}
	return s.FindMeta(ctx, key)
}

func (s *SignUpMetaStorage) FindMeta(ctx context.Context, signUpKey uuid.UUID) (*services.SignUpMeta, bool, error) {
	idKey := prefixSignUpMeta + signUpKey.String()
	metaResp := s.client.Get(ctx, idKey)
	if err := metaResp.Err(); err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("redis get sign up meta failed: %s", err)
	}

	meta := new(services.SignUpMeta)
	if err := json.Unmarshal([]byte(metaResp.Val()), meta); err != nil {
		return nil, false, fmt.Errorf("umarshalling sign up meta failed: %s", err)
	}

	return meta, true, nil
}

func (s *SignUpMetaStorage) Store(ctx context.Context, meta *services.SignUpMeta) error {
	metaJson, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("sign up metadata json marshalling failed: %s", err)
	}

	key := meta.SignUpKey.String()
	metaKey := prefixSignUpMeta + key
	phoneKey := prefixPhoneToSignUpKey + meta.Phone

	status := s.client.Set(ctx, metaKey, metaJson, s.conf.MetaLifetime)
	if err := status.Err(); err != nil {
		return err
	}

	status = s.client.Set(ctx, phoneKey, key, s.conf.MetaLifetime)
	if err := status.Err(); err != nil {
		return err
	}

	return nil
}
