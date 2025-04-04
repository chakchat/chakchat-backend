package instrumentation

import (
	"context"
	"errors"
	"strings"

	"github.com/chakchat/chakchat-backend/messaging-service/internal/application/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumName = "github.com/chakchat/chakchat-backend/messaging-service/internal/infrastructure/postgres/instrumentation"

	spanAttrStatement = "db.statement"
	spanAttrInsideTx  = "db.inside_tx"

	opSelect   = "SELECT"
	opInsert   = "INSERT"
	opUpdate   = "UPDATE"
	opDelete   = "DELETE"
	opAlter    = "ALTER"
	opDrop     = "DROP"
	opCreate   = "CREATE"
	opBegin    = "BEGIN"
	opCommit   = "COMMIT"
	opRollback = "ROLLBACK"

	opUnknown = "UNKNOWN_DB_QUERY"
)

func Tracing(db storage.SQLer) storage.SQLer {
	tracer := otel.Tracer(instrumName)

	return &tracingDB{
		db:     db,
		tracer: tracer,
	}
}

type tracingDB struct {
	db     storage.SQLer
	tracer trace.Tracer
}

func (r *tracingDB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	ctx, span := r.tracer.Start(ctx, getOpStmt(query), trace.WithAttributes(
		attribute.String(spanAttrStatement, query),
		attribute.Bool(spanAttrInsideTx, false),
	))
	defer span.End()

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
	}

	return res, err
}

func (r *tracingDB) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	ctx, span := r.tracer.Start(ctx, getOpStmt(query), trace.WithAttributes(
		attribute.String(spanAttrStatement, query),
		attribute.Bool(spanAttrInsideTx, false),
	))
	defer span.End()

	res, err := r.db.Query(ctx, query, args...)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
	}

	return res, err
}
func (r *tracingDB) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	ctx, span := r.tracer.Start(ctx, getOpStmt(query), trace.WithAttributes(
		attribute.String(spanAttrStatement, query),
		attribute.Bool(spanAttrInsideTx, false),
	))
	defer span.End()

	res := r.db.QueryRow(ctx, query, args...)
	// Here is a pitfall because the error is not recorded.

	return res
}

func (r *tracingDB) BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	ctx, span := r.tracer.Start(ctx, opBegin, trace.WithAttributes(
		attribute.Bool(spanAttrInsideTx, false),
	))
	defer span.End()

	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &tracingTx{
		tx:     tx,
		tracer: r.tracer,
	}, nil
}

func (r *tracingDB) Begin(ctx context.Context) (pgx.Tx, error) {
	ctx, span := r.tracer.Start(ctx, opBegin, trace.WithAttributes(
		attribute.Bool(spanAttrInsideTx, false),
	))
	defer span.End()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &tracingTx{
		tx:     tx,
		tracer: r.tracer,
	}, nil
}

type tracingTx struct {
	tx     pgx.Tx
	tracer trace.Tracer
}

func (t *tracingTx) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := t.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &tracingTx{
		tx:     tx,
		tracer: t.tracer,
	}, nil
}

func (t *tracingTx) Commit(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, opCommit, trace.WithAttributes(
		attribute.Bool(spanAttrInsideTx, true),
	))
	defer span.End()

	if err := t.tx.Commit(ctx); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (t *tracingTx) Rollback(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, opRollback, trace.WithAttributes(
		attribute.Bool(spanAttrInsideTx, true),
	))
	defer span.End()

	if err := t.tx.Rollback(ctx); err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

func (t *tracingTx) CopyFrom(
	ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource,
) (int64, error) {
	return t.tx.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (t *tracingTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return t.tx.SendBatch(ctx, b)
}

func (t *tracingTx) LargeObjects() pgx.LargeObjects {
	return t.tx.LargeObjects()
}

func (t *tracingTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return t.tx.Prepare(ctx, name, sql)
}

func (t *tracingTx) Exec(ctx context.Context, query string, args ...any) (commandTag pgconn.CommandTag, err error) {
	ctx, span := t.tracer.Start(ctx, getOpStmt(query), trace.WithAttributes(
		attribute.String(spanAttrStatement, query),
		attribute.Bool(spanAttrInsideTx, true),
	))
	defer span.End()

	res, err := t.tx.Exec(ctx, query, args...)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
	}

	return res, err
}

func (t *tracingTx) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	ctx, span := t.tracer.Start(ctx, getOpStmt(query), trace.WithAttributes(
		attribute.String(spanAttrStatement, query),
		attribute.Bool(spanAttrInsideTx, true),
	))
	defer span.End()

	res, err := t.tx.Query(ctx, query, args...)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
	}

	return res, err
}

func (t *tracingTx) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	ctx, span := t.tracer.Start(ctx, getOpStmt(query), trace.WithAttributes(
		attribute.String(spanAttrStatement, query),
		attribute.Bool(spanAttrInsideTx, true),
	))
	defer span.End()

	res := t.tx.QueryRow(ctx, query, args...)
	// Here is a pitfall because the error is not recorded.

	return res
}

func (t *tracingTx) Conn() *pgx.Conn { return t.tx.Conn() }

var opsStmt = [...]string{
	opSelect,
	opUpdate,
	opInsert,
	opDelete,
	opAlter,
	opDelete,
	opCreate,
	opBegin,
	opCommit,
	opRollback,
}

func getOpStmt(query string) string {
	firstPos, firstOp := uint(len(query)), opUnknown

	for _, op := range opsStmt {
		// uint() casting to make -1 be max uint value. Helps to ignore it in < comparison
		pos := uint(strings.Index(query, op))
		if pos < firstPos {
			firstPos, firstOp = pos, op
		}
	}

	return firstOp
}
