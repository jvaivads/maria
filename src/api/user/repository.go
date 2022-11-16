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
	SelectByID(int64) (User, error)
	//update(User) (User, error)
	//insert(User) (User, error)
}

func NewRelationalDB(client db.Client) Persister {
	return &relationalDB{
		client: client,
	}
}

type relationalDB struct {
	client db.Client
}

func (db *relationalDB) SelectByID(userID int64) (User, error) {
	var (
		u   User
		row *sql.Row
		err error

		query = `SELECT user_id FROM User WHERE id = ?`
	)

	if row = db.client.QueryRow(query, userID); row == nil {
		return u, emptyResultError(query)
	}

	if err = row.Scan(&u); err != nil {
		return u, scanError(err, query)
	}

	return u, nil
}
