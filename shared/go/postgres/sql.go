package postgres

import (
	"context"

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
