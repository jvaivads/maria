package db

import (
	"database/sql/driver"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

func SetClientQueryRowMock(
	rows *sqlmock.Rows,
	query string,
	scanError error,
	params ...driver.Value,
) func(m sqlmock.Sqlmock) func() error {
	return func(m sqlmock.Sqlmock) func() error {
		m.ExpectQuery(query).WithArgs(params...).WillReturnRows(rows)
		if scanError != nil {
			rows.RowError(0, scanError)
		}
		return func() error {
			return m.ExpectationsWereMet()
		}
	}
}

func SetClientQueryMock(
	rows *sqlmock.Rows,
	query string,
	queryError error,
	rowsError error,
	params ...driver.Value,
) func(m sqlmock.Sqlmock) func() error {
	return func(m sqlmock.Sqlmock) func() error {
		exp := m.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(params...)
		if queryError == nil {
			exp.WillReturnRows(rows)
			if rowsError != nil {
				rows.RowError(0, rowsError)
			}
		} else {
			exp.WillReturnError(queryError)
		}

		return func() error {
			return m.ExpectationsWereMet()
		}
	}
}

func SetClientExecMock(
	result driver.Result,
	query string,
	queryError error,
	params ...driver.Value,
) func(m sqlmock.Sqlmock) func() error {
	return func(m sqlmock.Sqlmock) func() error {
		exp := m.ExpectExec(regexp.QuoteMeta(query)).WithArgs(params...)
		if queryError == nil {
			exp.WillReturnResult(result)
		} else {
			exp.WillReturnError(queryError)
		}

		return func() error {
			return m.ExpectationsWereMet()
		}
	}
}
