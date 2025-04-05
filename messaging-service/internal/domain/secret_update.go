package domain

import (
	"time"
)

var TimeFunc = func() time.Time {
	return time.Now()
}

type (
	SecretKeyHash        []byte
	Encrypted            []byte
	InitializationVector []byte
)

type SecretData struct {
	KeyHash SecretKeyHash
	Payload Encrypted
	IV      InitializationVector
}

type SecretUpdate struct {
	Update
	Data SecretData
}

func (u *SecretUpdate) Delete(chat SecretChatter, sender UserID) error {
	if !chat.IsMember(sender) {
		return ErrUserNotMember
	}

	if chat.ChatID() != u.ChatID {
		return ErrUpdateNotFromChat
	}

	if u.DeletedFor(sender) {
		return ErrUpdateDeleted
	}

	u.AddDeletion(sender, DeleteModeForAll)
	return nil
}

func (u *SecretUpdate) Expired(exp *time.Duration) bool {
	if exp == nil {
		return false
	}
	now := TimeFunc()
	expTime := u.CreatedAt.Time().Add(*exp)
	return expTime.Before(now)
}

func NewSecretUpdate(chat SecretChatter, sender UserID, data SecretData) (*SecretUpdate, error) {
	if err := chat.ValidateCanSend(sender); err != nil {
		return nil, err
	}

	return &SecretUpdate{
		Update: Update{
			ChatID:   chat.ChatID(),
			SenderID: sender,
		},
		Data: data,
	}, nil
}
