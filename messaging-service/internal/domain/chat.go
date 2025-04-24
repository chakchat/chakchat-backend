package domain

import (
	"errors"

	"github.com/google/uuid"
)

const (
	ChatTypePersonal       = "personal"
	ChatTypeGroup          = "group"
	ChatTypeSecretPersonal = "secret_personal"
	ChatTypeSecretGroup    = "secret_group"
)

const (
	maxGroupNameLen   = 50
	maxDescriptionLen = 300
)

type (
	URL    string
	ChatID uuid.UUID
	UserID uuid.UUID
)

func NewUserID(id string) (UserID, error) {
	userId, err := uuid.Parse(id)
	return UserID(userId), err
}

func NewChatID() ChatID {
	return ChatID(uuid.New())
}

type Chat struct {
	ID        ChatID
	CreatedAt Timestamp
}

func (c *Chat) ChatID() ChatID {
	return c.ID
}

type Chatter interface {
	ChatID() ChatID
	IsMember(UserID) bool
	ValidateCanSend(UserID) error
}

func NormilizeMembers(members []UserID) []UserID {
	met := make(map[UserID]struct{}, len(members))
	normMembers := make([]UserID, 0, len(members))

	for _, member := range members {
		if _, ok := met[member]; !ok {
			normMembers = append(normMembers, member)
			met[member] = struct{}{}
		}
	}

	return normMembers
}

func ValidateGroupInfo(name, description string) error {
	var errs []error
	if name == "" {
		errs = append(errs, ErrGroupNameEmpty)
	}
	if len(name) > maxGroupNameLen {
		errs = append(errs, ErrGroupNameTooLong)
	}
	if len(description) > maxDescriptionLen {
		errs = append(errs, ErrGroupDescTooLong)
	}
	return errors.Join(errs...)
}
