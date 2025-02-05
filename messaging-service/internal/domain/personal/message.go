package personal

import "github.com/chakchat/chakchat-backend/messaging-service/internal/domain"

func (c *PersonalChat) validateCanSend(sender domain.UserID) error {
	if !c.IsMember(sender) {
		return domain.ErrUserNotMember
	}
	if c.Blocked() {
		return ErrChatBlocked
	}
	return nil
}

func (c *PersonalChat) validateCanReply(sender domain.UserID, replyTo *domain.Message) error {
	if replyTo.DeletedFor(sender) {
		return domain.ErrUpdateDeleted
	}
	if c.ChatID != replyTo.ChatID {
		return domain.ErrUpdateNotFromChat
	}
	return nil
}

func (c *PersonalChat) deleteMessage(sender domain.UserID, m *domain.Update, mode domain.DeleteMode) error {
	if err := c.validateCanSend(sender); err != nil {
		return err
	}

	if c.ChatID != m.ChatID {
		return domain.ErrUpdateNotFromChat
	}

	if m.DeletedFor(sender) {
		return domain.ErrUpdateDeleted
	}

	m.AddDeletion(sender, mode)
	return nil
}
