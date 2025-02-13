package secpersonal

import (
	"time"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type SecretPersonalChat struct {
	domain.SecretChat
	Members [2]domain.UserID
}

func NewSecretPersonalChatService(users [2]domain.UserID) (*SecretPersonalChat, error) {
	if users[0] == users[1] {
		return nil, domain.ErrChatWithMyself
	}

	return &SecretPersonalChat{
		SecretChat: domain.SecretChat{
			Chat: domain.Chat{
				ID: domain.NewChatID(),
			},
		},
		Members: users,
	}, nil
}

func (c *SecretPersonalChat) SetExpiration(sender domain.UserID, exp *time.Duration) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	c.Exp = exp
	return nil
}

func (c *SecretPersonalChat) Delete(sender domain.UserID) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	return nil
}

func (c *SecretPersonalChat) IsMember(user domain.UserID) bool {
	return user == c.Members[0] || user == c.Members[1]
}

func (c *SecretPersonalChat) ValidateCanSend(sender domain.UserID) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	return nil
}
