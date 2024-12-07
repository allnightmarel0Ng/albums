package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	Exec(ctx context.Context, sql string, args ...interface{}) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
}

type transaction struct {
	tx pgx.Tx
}

func (t *transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *transaction) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := t.tx.Exec(ctx, sql, args...)
	return err
}

func (t *transaction) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	row := t.tx.QueryRow(ctx, sql, args...)
	return NewRow(row)
}
