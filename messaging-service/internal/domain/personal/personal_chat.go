package personal

import (
	"errors"
	"slices"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

var (
	ErrAlreadyBlocked   = errors.New("chat is already blocked")
	ErrAlreadyUnblocked = errors.New("chat is already unblocked")

	ErrChatWithMyself = errors.New("chat with myself")
)

type PersonalChat struct {
	domain.Chat
	Members [2]domain.UserID

	Secret    bool
	BlockedBy []domain.UserID
}

func NewPersonalChat(users [2]domain.UserID) (*PersonalChat, error) {
	if users[0] == users[1] {
		return nil, ErrChatWithMyself
	}

	return &PersonalChat{
		Chat: domain.Chat{
			ChatID: domain.NewChatID(),
		},
		Members:   users,
		Secret:    false,
		BlockedBy: nil,
	}, nil
}

func (c *PersonalChat) BlockBy(user domain.UserID) error {
	if !c.IsMember(user) {
		return domain.ErrUserNotMember
	}

	if slices.Contains(c.BlockedBy, user) {
		return ErrAlreadyBlocked
	}

	c.BlockedBy = append(c.BlockedBy, user)

	return nil
}

func (c *PersonalChat) UnblockBy(user domain.UserID) error {
	if !c.IsMember(user) {
		return domain.ErrUserNotMember
	}

	if !slices.Contains(c.BlockedBy, user) {
		return ErrAlreadyUnblocked
	}

	c.BlockedBy = slices.DeleteFunc(c.BlockedBy, func(member domain.UserID) bool {
		return member == user
	})

	return nil
}

func (c *PersonalChat) Blocked() bool {
	return len(c.BlockedBy) > 0
}

func (c *PersonalChat) IsMember(user domain.UserID) bool {
	return user == c.Members[0] || user == c.Members[1]
}
