package domain

// If you want to add reaction to this list it says that it should be made more flexible
var allowedReactionTypes = map[ReactionType]struct{}{
	"heart":   {},
	"like":    {},
	"thunder": {},
	"cry":     {},
	"dislike": {},
	"bzZZ":    {},
}

type ReactionType string

type Reaction struct {
	Update

	Type      ReactionType
	MessageID UpdateID
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

	if err := validateReactionType(reaction); err != nil {
		return nil, err
	}

	return &Reaction{
		Update: Update{
			ChatID:   chat.ChatID(),
			SenderID: sender,
		},
		Type:      reaction,
		MessageID: m.UpdateID,
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

func validateReactionType(r ReactionType) error {
	if _, ok := allowedReactionTypes[r]; ok {
		return nil
	}
	return ErrInvalidReactionType
}
