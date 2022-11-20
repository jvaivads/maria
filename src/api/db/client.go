package db

import "database/sql"

type Client interface {
	Query(query string, params ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}
