package postgres

import (
	"context"
	"errors"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
)

type TxProviderStub struct{}

func (TxProviderStub) New(ctx context.Context) (storage.Tx, context.Context, error) {
	return TxStub{}, ctx, nil
}

type TxStub struct{}

func (TxStub) Commit(context.Context) error { return nil }
func (TxStub) Rollback(context.Context) error {
	return errors.New("tx stub: I cannot rollback - I am a stub")
}
