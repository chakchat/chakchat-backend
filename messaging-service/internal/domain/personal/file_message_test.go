package personal

import (
	"testing"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestFileMessage(t *testing.T) {
	user1, _ := domain.NewUserID("3d7ca3ef-3b0d-4113-91c9-20b7bf874324")
	user2, _ := domain.NewUserID("ce30ebc7-4058-4351-9a8f-66c71f987fdf")
	user3, _ := domain.NewUserID("fb048277-ad4f-4730-88eb-5e453c9ca5ce")
	chat := PersonalChat{
		Chat: domain.Chat{
			ChatID: domain.NewChatID(),
		},
		Members: [2]domain.UserID{user1, user2},
	}

	file := &FileMeta{
		FileId:   [16]byte{},
		FileName: "text.txt",
		MimeType: "plain/txt",
		FileSize: 1000,
		FileUrl:  "https://url.ru",
	}

	t.Run("New", func(t *testing.T) {
		_, err := chat.NewFileMessage(user3, file)
		require.ErrorIs(t, err, domain.ErrUserNotMember)

		msg1, err := chat.NewFileMessage(user1, file)
		require.NoError(t, err)
		require.Equal(t, chat.ChatID, msg1.ChatID)
		require.Equal(t, user1, msg1.SenderID)
	})

	t.Run("Delete", func(t *testing.T) {
		msg := FileMessage{
			Update: domain.Update{
				UpdateID: 12,
				ChatID:   chat.ChatID,
				SenderID: user1,
			},
			File: *file,
		}

		err := chat.DeleteFileMessage(user3, &msg, domain.DeleteModeForAll)
		require.ErrorIs(t, err, domain.ErrUserNotMember)

		err = chat.DeleteFileMessage(user2, &msg, domain.DeleteModeForSender)
		require.NoError(t, err)
		require.True(t, msg.DeletedFor(user2))
		require.False(t, msg.DeletedFor(user1))

		err = chat.DeleteFileMessage(user1, &msg, domain.DeleteModeForAll)
		require.NoError(t, err)
		require.True(t, msg.DeletedFor(user2))
	})
}
