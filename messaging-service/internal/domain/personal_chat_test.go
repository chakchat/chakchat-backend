package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPersonalChat_New_Success(t *testing.T) {
	user1 := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))
	user2 := UserID(uuid.MustParse("bdb22462-1a18-434e-8b68-91b997bf9553"))

	chat, err := NewPersonalChat([2]UserID{user1, user2})

	require.NoError(t, err)
	require.Equal(t, [2]UserID{user1, user2}, chat.Members)
	require.NotZero(t, chat.ID)
	require.False(t, chat.Secret)
	require.False(t, chat.Blocked)
}

func TestPersonalChat_New_Fails(t *testing.T) {
	user := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))

	_, err := NewPersonalChat([2]UserID{user, user})
	require.Equal(t, ErrChatWithMyself, err)
}

func TestPersonalChat_NewSecret_Success(t *testing.T) {
	user1 := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))
	user2 := UserID(uuid.MustParse("bdb22462-1a18-434e-8b68-91b997bf9553"))

	chat, err := NewSecretPersonalChat([2]UserID{user1, user2})

	require.NoError(t, err)
	require.Equal(t, [2]UserID{user1, user2}, chat.Members)
	require.NotZero(t, chat.ID)
	require.True(t, chat.Secret)
	require.False(t, chat.Blocked)
}

func TestPersonalChat_NewSecret_Fails(t *testing.T) {
	user := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))

	_, err := NewSecretPersonalChat([2]UserID{user, user})
	require.Equal(t, ErrChatWithMyself, err)
}

func TestPersonalChat_BlockBy_Success(t *testing.T) {
	user1 := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))
	user2 := UserID(uuid.MustParse("bdb22462-1a18-434e-8b68-91b997bf9553"))

	chat := PersonalChat{
		ID:        ChatID(uuid.Nil),
		Members:   [2]UserID{user1, user2},
		Blocked:   false,
		BlockedBy: []UserID{},
	}

	err := chat.BlockBy(user1)

	require.NoError(t, err)
	require.True(t, chat.Blocked)
	require.Equal(t, []UserID{user1}, chat.BlockedBy)
}

func TestPersonalChat_BlockBy_Fails(t *testing.T) {
	user1 := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))
	user2 := UserID(uuid.MustParse("bdb22462-1a18-434e-8b68-91b997bf9553"))
	user3 := UserID(uuid.MustParse("d8ffd2dd-451a-489a-ac8a-df84186004eb"))

	tests := []struct {
		Chat        PersonalChat
		UserID      UserID
		ExpectedErr error
	}{
		{
			Chat: PersonalChat{
				ID:        ChatID(uuid.Nil),
				Members:   [2]UserID{user1, user2},
				Blocked:   false,
				BlockedBy: []UserID{},
			},
			UserID:      user3,
			ExpectedErr: ErrUserNotMember,
		},
		{
			Chat: PersonalChat{
				ID:        ChatID(uuid.Nil),
				Members:   [2]UserID{user1, user2},
				Blocked:   true,
				BlockedBy: []UserID{user1},
			},
			UserID:      user1,
			ExpectedErr: ErrAlreadyBlocked,
		},
	}

	for _, test := range tests {
		err := test.Chat.BlockBy(test.UserID)
		require.ErrorIs(t, err, test.ExpectedErr)
	}
}

func TestPersonalChat_UnblockBy_Success(t *testing.T) {
	user1 := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))
	user2 := UserID(uuid.MustParse("bdb22462-1a18-434e-8b68-91b997bf9553"))

	chat := PersonalChat{
		ID:        ChatID(uuid.Nil),
		Members:   [2]UserID{user1, user2},
		Blocked:   true,
		BlockedBy: []UserID{user1, user2},
	}

	err := chat.UnblockBy(user1)

	require.NoError(t, err)
	require.True(t, chat.Blocked)
	require.Equal(t, []UserID{user2}, chat.BlockedBy)

	err = chat.UnblockBy(user2)

	require.NoError(t, err)
	require.False(t, chat.Blocked)
	require.Empty(t, chat.BlockedBy)
}

func TestPersonalChat_UnblockBy_Fails(t *testing.T) {
	user1 := UserID(uuid.MustParse("4f5d5d50-585b-4bd5-bb1c-82ca01427f8c"))
	user2 := UserID(uuid.MustParse("bdb22462-1a18-434e-8b68-91b997bf9553"))
	user3 := UserID(uuid.MustParse("d8ffd2dd-451a-489a-ac8a-df84186004eb"))

	tests := []struct {
		Chat        PersonalChat
		UserID      UserID
		ExpectedErr error
	}{
		{
			Chat: PersonalChat{
				ID:      ChatID(uuid.Nil),
				Members: [2]UserID{user1, user2},
			},
			UserID:      user3,
			ExpectedErr: ErrUserNotMember,
		},
		{
			Chat: PersonalChat{
				ID:        ChatID(uuid.Nil),
				Members:   [2]UserID{user1, user2},
				Blocked:   true,
				BlockedBy: []UserID{user1},
			},
			UserID:      user2,
			ExpectedErr: ErrAlreadyUnblocked,
		},
	}

	for _, test := range tests {
		err := test.Chat.UnblockBy(test.UserID)
		require.ErrorIs(t, err, test.ExpectedErr)
	}
}
