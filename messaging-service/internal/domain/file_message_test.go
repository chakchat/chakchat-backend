package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileMessage(t *testing.T) {
	user1, _ := NewUserID("3d7ca3ef-3b0d-4113-91c9-20b7bf874324")
	user2, _ := NewUserID("ce30ebc7-4058-4351-9a8f-66c71f987fdf")
	user3, _ := NewUserID("fb048277-ad4f-4730-88eb-5e453c9ca5ce")
	chat := &FakeChat{
		Chat: Chat{
			ID: NewChatID(),
		},
		Members: [2]UserID{user1, user2},
	}

	file := &FileMeta{
		FileId:   [16]byte{},
		FileName: "text.txt",
		MimeType: "plain/txt",
		FileSize: 1000,
		FileUrl:  "https://url.ru",
	}

	t.Run("New", func(t *testing.T) {
		_, err := NewFileMessage(chat, user3, file, nil)
		require.ErrorIs(t, err, ErrUserNotMember)

		msg1, err := NewFileMessage(chat, user1, file, nil)
		require.NoError(t, err)
		require.Equal(t, chat.ChatID(), msg1.ChatID)
		require.Equal(t, user1, msg1.SenderID)
	})

	t.Run("Delete", func(t *testing.T) {
		msg := FileMessage{
			Message: Message{
				Update: Update{
					UpdateID: 12,
					ChatID:   chat.ChatID(),
					SenderID: user1,
				},
			},
		}

		err := msg.Delete(chat, user3, DeleteModeForAll)
		require.ErrorIs(t, err, ErrUserNotMember)

		err = msg.Delete(chat, user2, DeleteModeForSender)
		require.NoError(t, err)
		require.True(t, msg.DeletedFor(user2))
		require.False(t, msg.DeletedFor(user1))

		err = msg.Delete(chat, user1, DeleteModeForAll)
		require.NoError(t, err)
		require.True(t, msg.DeletedFor(user2))
	})
}

type FakeChat struct {
	Chat
	Members [2]UserID
}

func (c *FakeChat) ChatID() ChatID {
	return c.Chat.ID
}

func (c *FakeChat) ValidateCanSend(user UserID) error {
	if !c.IsMember(user) {
		return ErrUserNotMember
	}
	return nil
}

func (c *FakeChat) IsMember(user UserID) bool {
	return user == c.Members[0] || user == c.Members[1]
}
