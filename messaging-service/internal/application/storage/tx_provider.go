package storage

import "context"

type TxProvider interface {
	New(context.Context) (Tx, context.Context, error)
}

type Tx interface {
	Commit() error
	Rollback() error
}
