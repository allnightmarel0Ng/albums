package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Database interface {
	Query(sql string, args ...interface{}) (Rows, error)
	QueryRow(sql string, args ...interface{}) Row

	Close()

	Begin() error
	Commit() error
	Rollback() error
}

type db struct {
	pool *pgxpool.Pool
	ctx  context.Context
	tx   pgx.Tx
}

func NewDatabase(ctx context.Context, connectionString string) (Database, error) {
	pool, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	return &db{
		pool: pool,
		ctx:  ctx,
		tx:   nil,
	}, nil
}

func (db *db) Query(sql string, args ...interface{}) (Rows, error) {
	var rows pgx.Rows
	var err error

	switch {
	case db.tx == nil:
		rows, err = db.pool.Query(db.ctx, sql, args...)
	case db.tx != nil:
		rows, err = db.tx.Query(db.ctx, sql, args...)
	}

	if err != nil {
		return nil, err
	}

	return NewRows(rows), nil
}

func (db *db) QueryRow(sql string, args ...interface{}) Row {
	var row pgx.Row

	switch {
	case db.tx == nil:
		row = db.pool.QueryRow(db.ctx, sql, args...)
	case db.tx != nil:
		row = db.pool.QueryRow(db.ctx, sql, args...)
	}

	return NewRow(row)
}

func (db *db) Close() {
	db.pool.Close()
}

func (db *db) Begin() error {
	tx, err := db.pool.Begin(db.ctx)
	if err != nil {
		return err
	}

	db.tx = tx
	return nil
}

func (db *db) Commit() error {
	err := db.tx.Commit(db.ctx)
	if err != nil {
		return err
	}

	db.tx = nil
	return nil
}

func (db *db) Rollback() error {
	err := db.tx.Rollback(db.ctx)
	if err != nil {
		return err
	}

	db.tx = nil
	return nil
}
