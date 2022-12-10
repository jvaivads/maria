package db

import (
	"database/sql"
	"fmt"
)

// Next functions are wrapper of sql package's errors, they are used for keeping query
var (
	QueryError = func(err error, query string) error {
		return fmt.Errorf("unexpected error querying result by using query: '%s'. error: %w", query, err)
	}
	RowsError = func(err error, query string) error {
		return fmt.Errorf("unexpected rows error by using query: '%s'. error: %w", query, err)
	}
	ExecError = func(err error, query string) error {
		return fmt.Errorf("unexpected exec error by using query: '%s'. error: %w", query, err)
	}
	LastInsertedError = func(err error, query string) error {
		return fmt.Errorf("unexpected last inserted error by using query: '%s'. error: %w", query, err)
	}
	RowsAffectedError = func(err error, query string) error {
		return fmt.Errorf("unexpected getting rows affected error by using query: '%s'. error: %w", query, err)
	}
	ScanError = func(err error, query string) error {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("unexpected error scannig result by using query: '%s'. error: %w", query, err)
	}
)
