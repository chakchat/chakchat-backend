package personal

import (
	"testing"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestReaction(t *testing.T) {
	user1, _ := domain.NewUserID("3d7ca3ef-3b0d-4113-91c9-20b7bf874324")
	user2, _ := domain.NewUserID("ce30ebc7-4058-4351-9a8f-66c71f987fdf")
	user3, _ := domain.NewUserID("fb048277-ad4f-4730-88eb-5e453c9ca5ce")
	chat := PersonalChat{
		Chat: domain.Chat{
			ID: domain.NewChatID(),
		},
		Members: [2]domain.UserID{user1, user2},
	}

	txtMsg := domain.TextMessage{
		Message: domain.Message{
			Update: domain.Update{
				UpdateID: 12,
				ChatID:   chat.ID,
				SenderID: user1,
			},
		},
	}

	_, err := chat.NewReaction(user3, &txtMsg.Message, "some_reaction_idk")
	require.ErrorIs(t, err, domain.ErrUserNotMember)

	reaction, err := chat.NewReaction(user1, &txtMsg.Message, "some_reaction_idk")
	require.NoError(t, err)
	require.Equal(t, chat.ID, reaction.ChatID)
	require.Equal(t, user1, reaction.SenderID)

	err = chat.DeleteReaction(user3, &reaction)
	require.ErrorIs(t, err, domain.ErrUserNotMember)

	err = chat.DeleteReaction(user2, &reaction)
	require.ErrorIs(t, err, ErrReactionNotFromUser)

	err = chat.DeleteReaction(user1, &reaction)
	require.NoError(t, err)
	require.True(t, reaction.DeletedFor(user1))
	require.True(t, reaction.DeletedFor(user2))
}
