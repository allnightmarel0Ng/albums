package postgres

import "github.com/jackc/pgx/v4"

type Row interface {
	Scan(dest ...interface{}) error
}

type row struct {
	row pgx.Row
}

func NewRow(r pgx.Row) Row {
	return &row{
		row: r,
	}
}

func (r *row) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}
