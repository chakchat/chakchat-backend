package personal

import (
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type ReactionType string

var (
	ErrReactionNotFromUser = errors.New("the reaction is not from this user")
)

type Reaction struct {
	domain.Update

	Type ReactionType
}

func (c *PersonalChat) NewReaction(
	sender domain.UserID,
	m *domain.Message,
	reaction ReactionType,
) (Reaction, error) {
	if err := c.validateCanSend(sender); err != nil {
		return Reaction{}, err
	}

	if c.ChatID != m.ChatID {
		return Reaction{}, domain.ErrUpdateNotFromChat
	}

	if m.DeletedFor(sender) {
		return Reaction{}, domain.ErrUpdateDeleted
	}

	return Reaction{
		Update: domain.Update{
			ChatID:   c.ChatID,
			SenderID: sender,
		},
		Type: reaction,
	}, nil
}

func (c *PersonalChat) DeleteReaction(sender domain.UserID, r *Reaction) error {
	if err := c.validateCanSend(sender); err != nil {
		return err
	}

	if c.ChatID != r.ChatID {
		return domain.ErrUpdateNotFromChat
	}

	if r.SenderID != sender {
		return ErrReactionNotFromUser
	}

	if r.DeletedFor(sender) {
		return domain.ErrUpdateDeleted
	}

	r.AddDeletion(sender, domain.DeleteModeForAll)
	return nil
}
