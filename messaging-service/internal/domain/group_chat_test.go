package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGroupChat_New(t *testing.T) {
	user1 := UserID(uuid.MustParse("601ff05e-35cb-43dd-83b4-2d2d6d2f5773"))
	user2 := UserID(uuid.MustParse("b0c54326-219c-4865-91b9-bdb8ac97e85c"))

	t.Run("Success", func(t *testing.T) {
		group, err := NewGroupChat(
			user1,
			[]UserID{user1, user1, user2, user2},
			"group name")

		require.NoError(t, err)
		require.Equal(t, []UserID{user1, user2}, group.Members)
		require.Equal(t, user1, group.Admin)
	})

	t.Run("AdminNotMember", func(t *testing.T) {
		_, err := NewGroupChat(
			user1,
			[]UserID{user2},
			"group name")

		require.ErrorIs(t, err, ErrAdminNotMember)
	})

	t.Run("InvalidName", func(t *testing.T) {
		_, err := NewGroupChat(
			user1,
			[]UserID{user1, user2},
			"",
		)

		require.ErrorIs(t, err, ErrGroupNameEmpty)
	})
}
func TestGroupChat_UpdateInfo_FailsValidation(t *testing.T) {
	g := GroupChat{
		ID:      ChatID(uuid.New()),
		Admin:   UserID(uuid.MustParse("4b6e828d-7c6a-4bf2-8a2c-19f13ece23cc")),
		Members: []UserID{UserID(uuid.MustParse("4b6e828d-7c6a-4bf2-8a2c-19f13ece23cc"))},
		Secret:  false,
		Name:    "Super puper group",
	}

	tests := []struct {
		Name, Desc string
		TargetErr  error
	}{
		{
			Name:      "",
			Desc:      "",
			TargetErr: ErrGroupNameEmpty,
		},
		{
			Name:      string(make([]byte, 51)),
			Desc:      "",
			TargetErr: ErrGroupNameTooLong,
		},
		{
			Name:      "new valid name",
			Desc:      string(make([]byte, 301)),
			TargetErr: ErrGroupDescTooLong,
		},
	}

	for _, test := range tests {
		err := g.UpdateInfo(test.Name, test.Desc)

		require.ErrorIs(t, err, test.TargetErr)
	}
}
