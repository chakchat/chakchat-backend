package repository

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/generic"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/chakchat/chakchat-backend/messaging-service/internal/domain"
)

type GenericUpdateRepository interface {
	// Should return ErrNotFound if not found
	GetLastUpdateID(
		context.Context,
		storage.ExecQuerier,
		domain.ChatID,
	) (domain.UpdateID, error)
	GetRange(
		ctx context.Context,
		db storage.ExecQuerier,
		visibleTo domain.UserID,
		chatID domain.ChatID,
		from, to domain.UpdateID,
	) ([]generic.Update, error)
	Get(
		ctx context.Context,
		db storage.ExecQuerier,
		visibleTo domain.UserID,
		chatID domain.ChatID,
		updateID domain.UpdateID,
	) (*generic.Update, error)
	FetchLast(
		ctx context.Context,
		db storage.ExecQuerier,
		visibleTo domain.UserID,
		chatID domain.ChatID,
		opts ...FetchLastOption,
	) ([]generic.Update, error)
}

// Specifies what updates are counted
//
// For example if [FetchLastModeMessagesReactions] is chosen and Count = 5
// then last updates will be fetched with 5 messages or reactions.
//
// Other update types are not counted but are also fetched.
type FetchLastMode int

const (
	// Count only messages
	FetchLastModeMessages FetchLastMode = iota
	// Count only messages and reactions
	FetchLastModeMessagesReactions
	// Count all update types
	FetchLastModeAll
)

type FetchLastOptions struct {
	Mode  FetchLastMode
	Count int
}

func NewFetchLastOptions(opts ...FetchLastOption) *FetchLastOptions {
	res := defaultFetchLastOptions // copy
	for _, optFunc := range opts {
		optFunc(&res)
	}
	return &res
}

var defaultFetchLastOptions = FetchLastOptions{
	Mode:  FetchLastModeMessages,
	Count: 5,
}

type FetchLastOption func(*FetchLastOptions)

// Defaults to 5
func WithFetchLastCount(count int) FetchLastOption {
	return func(opts *FetchLastOptions) {
		opts.Count = count
	}
}

// Defaults to FetchLastModeMessages
func WithFetchLastOptions(mode FetchLastMode) FetchLastOption {
	return func(opts *FetchLastOptions) {
		opts.Mode = mode
	}
}
