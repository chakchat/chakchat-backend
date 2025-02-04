package domain

import (
	"errors"
	"slices"
)

var (
	ErrAlreadyBlocked   = errors.New("chat is already blocked")
	ErrAlreadyUnblocked = errors.New("chat is already unblocked")
	ErrUserNotMember    = errors.New("user is not member of a chat")

	ErrChatWithMyself = errors.New("chat with myself")
)

type PersonalChat struct {
	Chat
	Members [2]UserID

	Secret    bool
	Blocked   bool
	BlockedBy []UserID
}

func NewPersonalChat(users [2]UserID) (*PersonalChat, error) {
	if users[0] == users[1] {
		return nil, ErrChatWithMyself
	}

	return &PersonalChat{
		Chat: Chat{
			ChatID: NewChatID(),
		},
		Members:   users,
		Secret:    false,
		Blocked:   false,
		BlockedBy: nil,
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

	if c.Blocked && len(c.BlockedBy) == 0 {
		c.Blocked = false
	}
	return nil
}

func (c *PersonalChat) IsMember(user UserID) bool {
	return user == c.Members[0] || user == c.Members[1]
}
