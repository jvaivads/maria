package user

import (
	"database/sql"
	"errors"
	"fmt"
	"maria/src/api/db"
)

var (
	emptyResultError = func(query string) error {
		return errors.New(fmt.Sprintf("db client returns nothing by using query: '%s'", query))
	}

	scanError = func(err error, query string) error {
		if err == sql.ErrNoRows {
			return fmt.Errorf("result is empty by using query: '%s'. error: %w", query, err)
		}

		return fmt.Errorf("unexpected error scannig result by using query: '%s'. error: %w", query, err)
	}
)

type Persister interface {
	selectByID(int64) (user, error)
	//update(user) (user, error)
	//insert(user) (user, error)
}

func newRelationalDB() *relationalDB {
	return &relationalDB{
		client: db.GetSQLClient(),
	}
}

type relationalDB struct {
	client db.Client
}

func (db *relationalDB) selectByID(userID int64) (user, error) {
	var (
		u   user
		row *sql.Row
		err error

		query = `SELECT user_id FROM user WHERE id = ?`
	)

	if row = db.client.QueryRow(query, userID); row == nil {
		return u, emptyResultError(query)
	}

	if err = row.Scan(&u); err != nil {
		return u, scanError(err, query)
	}

	return u, nil
}
