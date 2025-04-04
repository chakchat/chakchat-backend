package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SQLer interface {
	TxProvider
	ExecQuerier
}

type TxProvider interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Begin(context.Context) (pgx.Tx, error)
}

type ExecQuerier interface {
	Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(_ context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(_ context.Context, query string, args ...any) pgx.Row
}

// Should be used only in `defer` block because it has recover() call
func FinishTx(ctx context.Context, tx pgx.Tx, err *error) {
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
