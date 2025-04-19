package request

import (
	"time"

	"github.com/google/uuid"
)

type BlockChat struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
}

type UnblockChat struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
}

type CreatePersonalChat struct {
	SenderID uuid.UUID
	MemberID uuid.UUID
}

type CreateSecretPersonalChat struct {
	SenderID uuid.UUID
	MemberID uuid.UUID
}

type CreateGroup struct {
	SenderID uuid.UUID
	Members  []uuid.UUID
	Name     string
}

type SetExpiration struct {
	ChatID     uuid.UUID
	SenderID   uuid.UUID
	Expiration *time.Duration
}

type DeleteChat struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
}

type AddMember struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	MemberID uuid.UUID
}

type DeleteMember struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	MemberID uuid.UUID
}

type UpdateGroupPhoto struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
	FileID   uuid.UUID
}

type DeleteGroupPhoto struct {
	ChatID   uuid.UUID
	SenderID uuid.UUID
}

type UpdateGroupInfo struct {
	ChatID      uuid.UUID
	SenderID    uuid.UUID
	Name        string
	Description string
}

type CreateSecretGroup struct {
	SenderID uuid.UUID
	Members  []uuid.UUID
	Name     string
}

type UpdateSecretGroupInfo struct {
	ChatID      uuid.UUID
	SenderID    uuid.UUID
	Name        string
	Description string
}

// Options for getting chat
type GetChatOptions struct {
	LoadPreviewCount int
	LoadLastUpdateID bool
}

func NewGetChatOptions(opts ...GetChatOption) *GetChatOptions {
	res := new(GetChatOptions)
	for _, optFunc := range opts {
		optFunc(res)
	}

	return res
}

type GetChatOption func(*GetChatOptions)

// It will fetch last updates with repository.FetchLastModeMessages option
func WithChatPreview(count int) GetChatOption {
	return func(opts *GetChatOptions) {
		opts.LoadPreviewCount = count
	}
}

func WithChatLastUpdateID() GetChatOption {
	return func(opts *GetChatOptions) {
		opts.LoadLastUpdateID = true
	}
}
