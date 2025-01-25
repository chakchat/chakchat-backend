package domain

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAlreadyBlocked   = errors.New("chat is already blocked")
	ErrAlreadyUnblocked = errors.New("chat is already unblocked")
	ErrUserNotMember    = errors.New("user is not member of a chat")
)

type (
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Timestamp int64
	URL       string
)

var TimeFunc = func() time.Time {
	return time.Now()
}

type PersonalChat struct {
	ID      ChatID
	Members [2]UserID

	Secret    bool
	Blocked   bool
	BlockedBy []UserID
	CreatedAt Timestamp
}

func NewPersonalChat(users [2]UserID) (*PersonalChat, error) {
	return &PersonalChat{
		ID:        ChatID(uuid.New()),
		Members:   users,
		Secret:    false,
		Blocked:   false,
		BlockedBy: nil,
		CreatedAt: Timestamp(TimeFunc().Unix()),
	}, nil
}

func (c *PersonalChat) BlockBy(user UserID) error {
	if !c.IsMember(user) {
		return ErrUserNotMember
	}

	if slices.Contains(c.BlockedBy, user) {
		return ErrAlreadyBlocked
	}

	c.BlockedBy = append(c.BlockedBy, user)

	// This condition is only for producing an event in the future
	if !c.Blocked {
		c.Blocked = true
	}
	return nil
}

func (c *PersonalChat) UnblockBy(user UserID) error {
	if !c.IsMember(user) {
		return ErrUserNotMember
	}

	if !slices.Contains(c.BlockedBy, user) {
		return ErrAlreadyUnblocked
	}

	c.BlockedBy = slices.DeleteFunc(c.BlockedBy, func(member UserID) bool {
		return member == user
	})

	if c.Blocked {
		c.Blocked = false
	}
	return nil
}

func (c *PersonalChat) IsMember(user UserID) bool {
	return user == c.Members[0] || user == c.Members[1]
}
