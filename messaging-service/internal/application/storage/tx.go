package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TxProvider interface {
	BeginTx(context.Context) (Tx, error)
}

type ExecQuerier interface {
	Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(_ context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(q_ context.Context, uery string, args ...any) pgx.Row
}

type Tx interface {
	ExecQuerier
	Commit(context.Context) error
	Rollback(context.Context) error
}

// Should be used only in `defer` block because it has recover() call
func FinishTx(ctx context.Context, tx Tx, err *error) {
	if p := recover(); p != nil {
		if *err != nil {
			*err = fmt.Errorf(`[PANIC] panic recovered: %+v. (error overwritten: "%w")`, p, *err)
		} else {
			*err = fmt.Errorf(`[PANIC] panic recovered: %+v`, p)
		}
	}

	if *err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}
