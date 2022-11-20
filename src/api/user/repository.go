package user

import (
	"database/sql"
	"fmt"
	"maria/src/api/db"
)

const (
	getUserByIDQuery = `SELECT user_id FROM User WHERE id = ?`
)

var (
	scanError = func(err error, query string) error {
		if err == sql.ErrNoRows {
			return nil
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
	)

	row = db.client.QueryRow(getUserByIDQuery, userID)

	if err = row.Scan(&u.ID); err != nil {
		return u, scanError(err, getUserByIDQuery)
	}

	return u, nil
}
