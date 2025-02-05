package secpersonal

import "github.com/chakchat/chakchat-backend/messaging-service/internal/domain"

type SecretPersonalChat struct {
	domain.SecretChat
	Members [2]domain.UserID
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
