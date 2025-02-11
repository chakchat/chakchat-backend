package storage

import (
	"context"
)

type TxProvider interface {
	New(context.Context) (Tx, context.Context, error)
}

type Tx interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

func RunInTx(ctx context.Context, provider TxProvider, f func(context.Context) error) error {
	tx, ctx, err := provider.New(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback(ctx)
		}
		HandleTx(ctx, tx, err)
	}()
	return f(ctx)
}

func HandleTx(ctx context.Context, tx Tx, err error) {
	// TODO: handle tx errors.
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}
