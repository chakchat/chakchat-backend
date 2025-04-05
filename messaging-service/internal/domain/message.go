package domain

type Message struct {
	Update

	ReplyTo   *UpdateID
	Forwarded bool
}

func (m *Message) Delete(chat Chatter, sender UserID, mode DeleteMode) error {
	if err := chat.ValidateCanSend(sender); err != nil {
		return err
	}

	if chat.ChatID() != m.ChatID {
		return ErrUpdateNotFromChat
	}

	if m.DeletedFor(sender) {
		return ErrUpdateDeleted
	}

	m.AddDeletion(sender, mode)
	return nil
}

func validateCanReply(chat Chatter, sender UserID, replyTo *Message) error {
	if replyTo.DeletedFor(sender) {
		return ErrUpdateDeleted
	}
	if chat.ChatID() != replyTo.ChatID {
		return ErrUpdateNotFromChat
	}
	return nil
}
