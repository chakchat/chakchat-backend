package personal

import (
	"slices"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type PersonalChat struct {
	domain.Chat
	Members [2]domain.UserID

	BlockedBy []domain.UserID
}

func NewPersonalChat(users [2]domain.UserID) (*PersonalChat, error) {
	if users[0] == users[1] {
		return nil, domain.ErrChatWithMyself
	}

	return &PersonalChat{
		Chat: domain.Chat{
			ID: domain.NewChatID(),
		},
		Members:   users,
		BlockedBy: nil,
	}, nil
}

func (c *PersonalChat) Delete(sender domain.UserID) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	return nil
}

func (c *PersonalChat) BlockBy(user domain.UserID) error {
	if !c.IsMember(user) {
		return domain.ErrUserNotMember
	}

	if slices.Contains(c.BlockedBy, user) {
		return domain.ErrAlreadyBlocked
	}

	c.BlockedBy = append(c.BlockedBy, user)

	return nil
}

func (c *PersonalChat) UnblockBy(user domain.UserID) error {
	if !c.IsMember(user) {
		return domain.ErrUserNotMember
	}

	if !slices.Contains(c.BlockedBy, user) {
		return domain.ErrAlreadyUnblocked
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

func (c *PersonalChat) ValidateCanSend(sender domain.UserID) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	if c.Blocked() {
		return domain.ErrChatBlocked
	}
	return nil
}
