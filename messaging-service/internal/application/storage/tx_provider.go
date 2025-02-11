package storage

import "context"

type TxProvider interface {
	New(context.Context) (Tx, context.Context, error)
}

type Tx interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

func HandleTx(ctx context.Context, tx Tx, err error) {
	// TODO: handle tx errors.
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}
