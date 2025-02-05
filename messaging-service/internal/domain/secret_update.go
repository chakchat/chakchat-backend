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

type SecretData struct {
	KeyID   SecretKeyID
	Payload Encrypted
	IV      InitializationVector
}

type SecretUpdate struct {
	Update
	Data SecretData
}

func (u *SecretUpdate) Expired(exp time.Duration) bool {
	now := TimeFunc()
	expTime := u.CreatedAt.Time().Add(exp)
	return expTime.Before(now)
}

func NewSecretUpdate(chat SecretChatter, sender UserID, data SecretData) (SecretUpdate, error) {
	if err := chat.ValidateCanSend(sender); err != nil {
		return SecretUpdate{}, err
	}

	return SecretUpdate{
		Update: Update{
			ChatID:   chat.ChatID(),
			SenderID: sender,
		},
		Data: data,
	}, nil
}
