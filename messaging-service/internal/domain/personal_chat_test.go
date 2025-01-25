package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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
