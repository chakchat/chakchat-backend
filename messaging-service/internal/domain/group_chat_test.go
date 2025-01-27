package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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
