package personal

import (
	"testing"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTextMessage(t *testing.T) {
	user1, _ := domain.NewUserID("3d7ca3ef-3b0d-4113-91c9-20b7bf874324")
	user2, _ := domain.NewUserID("ce30ebc7-4058-4351-9a8f-66c71f987fdf")
	user3, _ := domain.NewUserID("fb048277-ad4f-4730-88eb-5e453c9ca5ce")
	chat := PersonalChat{
		Chat: domain.Chat{
			ChatID: domain.NewChatID(),
		},
		Members: [2]domain.UserID{user1, user2},
	}

	t.Run("New", func(t *testing.T) {
		_, err := chat.NewTextMessage(user1, "", nil)
		require.ErrorIs(t, err, ErrTextEmpty)

		_, err = chat.NewTextMessage(user3, "valid but not a member", nil)
		require.ErrorIs(t, err, domain.ErrUserNotMember)

		_, err = chat.NewTextMessage(user1, string(make([]byte, 3000)), nil)
		require.ErrorIs(t, err, ErrTooMuchTextRunes)

		msg1, err := chat.NewTextMessage(user1, "valid text message", nil)
		require.NoError(t, err)
		require.Equal(t, chat.ChatID, msg1.ChatID)
		require.Equal(t, user1, msg1.SenderID)
	})

	msgBase := TextMessage{
		Message: Message{
			Update: domain.Update{
				UpdateID: 12,
				ChatID:   chat.ChatID,
				SenderID: user1,
			},
		},
		Text: "previous text",
	}

	t.Run("Edit", func(t *testing.T) {
		msg := msgBase

		err := chat.EditTextMessage(user2, &msg, "valid but user is not a sender")
		require.ErrorIs(t, err, domain.ErrUserNotSender)

		err = chat.EditTextMessage(user3, &msg, "valid but not a member")
		require.ErrorIs(t, err, domain.ErrUserNotMember)
	})

	t.Run("EditDeleted", func(t *testing.T) {
		msg := msgBase
		msg.Deleted = []domain.UpdateDeleted{{
			Update: domain.Update{
				UpdateID: 13,
				ChatID:   chat.ChatID,
				SenderID: user1,
			},
			DeletedID: 12,
			Mode:      domain.DeleteModeForSender,
		}}

		err := chat.EditTextMessage(user1, &msg, "new valid text")
		require.ErrorIs(t, err, domain.ErrUpdateDeleted)
	})

	t.Run("Delete", func(t *testing.T) {
		msg := msgBase

		err := chat.DeleteTextMessage(user3, &msg, domain.DeleteModeForAll)
		require.ErrorIs(t, err, domain.ErrUserNotMember)

		err = chat.DeleteTextMessage(user2, &msg, domain.DeleteModeForSender)
		require.NoError(t, err)
		require.True(t, msg.DeletedFor(user2))
		require.False(t, msg.DeletedFor(user1))

		err = chat.DeleteTextMessage(user1, &msg, domain.DeleteModeForAll)
		require.NoError(t, err)
		require.True(t, msg.DeletedFor(user2))
	})
}
