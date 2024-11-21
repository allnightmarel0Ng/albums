package postgres

import "github.com/jackc/pgx/v4"

type Rows interface {
	Close()
	Next() bool
	Scan(dest ...interface{}) error
}

type rows struct {
	rows pgx.Rows
}

func NewRows(r pgx.Rows) Rows {
	return &rows{
		rows: r,
	}
}

func (r *rows) Close() {
	r.rows.Close()
}

func (r *rows) Next() bool {
	return r.rows.Next()
}

func (r *rows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}
