package secpersonal

import "github.com/chakchat/chakchat-backend/messaging-service/internal/domain"

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

func (c *SecretPersonalChat) IsMember(user domain.UserID) bool {
	return user == c.Members[0] || user == c.Members[1]
}

func (c *SecretPersonalChat) ValidateCanSend(sender domain.UserID) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	return nil
}
