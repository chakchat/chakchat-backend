package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaction(t *testing.T) {
	user1, _ := NewUserID("3d7ca3ef-3b0d-4113-91c9-20b7bf874324")
	user2, _ := NewUserID("ce30ebc7-4058-4351-9a8f-66c71f987fdf")
	user3, _ := NewUserID("fb048277-ad4f-4730-88eb-5e453c9ca5ce")
	chat := &FakeChat{
		Chat: Chat{
			ID: NewChatID(),
		},
		Members: [2]UserID{user1, user2},
	}

	txtMsg := TextMessage{
		Message: Message{
			Update: Update{
				UpdateID: 12,
				ChatID:   chat.ID,
				SenderID: user1,
			},
		},
	}

	_, err := NewReaction(chat, user3, &txtMsg.Message, "some_reaction_idk")
	require.ErrorIs(t, err, ErrUserNotMember)

	reaction, err := NewReaction(chat, user1, &txtMsg.Message, "some_reaction_idk")
	require.NoError(t, err)
	require.Equal(t, chat.ID, reaction.ChatID)
	require.Equal(t, user1, reaction.SenderID)

	err = reaction.Delete(chat, user3)
	require.ErrorIs(t, err, ErrUserNotMember)

	err = reaction.Delete(chat, user2)
	require.ErrorIs(t, err, ErrReactionNotFromUser)

	err = reaction.Delete(chat, user1)
	require.NoError(t, err)
	require.True(t, reaction.DeletedFor(user1))
	require.True(t, reaction.DeletedFor(user2))
}
