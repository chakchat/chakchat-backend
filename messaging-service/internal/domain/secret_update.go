package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	SecretKeyID          uuid.UUID
	Encrypted            []byte
	InitializationVector []byte
)

type SecretUpdate struct {
	Update

	KeyID   SecretKeyID
	Payload Encrypted
	IV      InitializationVector
}

func (u *SecretUpdate) Expired(exp time.Duration) bool {
	now := TimeFunc()
	expTime := u.CreatedAt.Time().Add(exp)
	return expTime.Before(now)
}
