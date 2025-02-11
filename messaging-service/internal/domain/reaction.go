package domain

type ReactionType string

type Reaction struct {
	Update

	Type ReactionType
}

func NewReaction(
	chat Chatter,
	sender UserID,
	m *Message,
	reaction ReactionType,
) (*Reaction, error) {
	if err := chat.ValidateCanSend(sender); err != nil {
		return nil, err
	}

	if chat.ChatID() != m.ChatID {
		return nil, ErrUpdateNotFromChat
	}

	if m.DeletedFor(sender) {
		return nil, ErrUpdateDeleted
	}

	return &Reaction{
		Update: Update{
			ChatID:   chat.ChatID(),
			SenderID: sender,
		},
		Type: reaction,
	}, nil
}

func (r *Reaction) Delete(chat Chatter, sender UserID) error {
	if err := chat.ValidateCanSend(sender); err != nil {
		return err
	}

	if chat.ChatID() != r.ChatID {
		return ErrUpdateNotFromChat
	}

	if r.SenderID != sender {
		return ErrReactionNotFromUser
	}

	if r.DeletedFor(sender) {
		return ErrUpdateDeleted
	}

	r.AddDeletion(sender, DeleteModeForAll)
	return nil
}
