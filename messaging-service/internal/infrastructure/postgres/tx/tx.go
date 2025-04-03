package tx

import (
	"context"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxProvider struct {
	db *pgxpool.Pool
}

func NewTxProvider(db *pgxpool.Pool) *TxProvider {
	return &TxProvider{db: db}
}

func (p *TxProvider) BeginTx(ctx context.Context) (storage.Tx, error) {
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})

	res := &pgTx{tx}

	return res, err
}

type pgTx struct {
	tx pgx.Tx
}

func (t *pgTx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgTx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *pgTx) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return t.tx.Exec(ctx, query, args...)
}

func (t *pgTx) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return t.tx.Query(ctx, query, args...)
}

func (t *pgTx) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return t.tx.QueryRow(ctx, query, args...)
}
