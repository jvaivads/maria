package user

import (
	"database/sql"
	"fmt"
	"maria/src/api/db"
	"time"
)

const (
	getUserByIDQuery  = `SELECT user_id, user_name, alias, email, active, date_created FROM user WHERE id = ?`
	getUserByAnyQuery = `SELECT user_id, user_name, alias, email, active, date_created FROM user WHERE user_name = ? OR alias = ? OR email = ?`
	insertUserQuery   = `INSERT INTO user (user_name, alias, email, active, date_created) VALUES (?, ?, ?, false, ?)`
)

var (
	queryError = func(err error, query string) error {
		return fmt.Errorf("unexpected error querying result by using query: '%s'. error: %w", query, err)
	}
	rowsError = func(err error, query string) error {
		return fmt.Errorf("unexpected rows error by using query: '%s'. error: %w", query, err)
	}
	execError = func(err error, query string) error {
		return fmt.Errorf("unexpected exec error by using query: '%s'. error: %w", query, err)
	}
	lastInsertedError = func(err error, query string) error {
		return fmt.Errorf("unexpected getting last inserted error by using query: '%s'. error: %w", query, err)
	}
	scanError = func(err error, query string) error {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("unexpected error scannig result by using query: '%s'. error: %w", query, err)
	}
)

type Persister interface {
	SelectByID(int64) (User, error)
	SelectByAny(string, string, string) ([]User, error)
	CreateUser(request NewUserRequest) (User, error)
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

	if err = row.Scan(
		&u.ID,
		&u.UserName,
		&u.Alias,
		&u.Email,
		&u.Active,
		&u.DateCreated,
	); err != nil {
		return u, scanError(err, getUserByIDQuery)
	}

	return u, nil
}

func (db *relationalDB) SelectByAny(name, alias, email string) ([]User, error) {
	var (
		rows  *sql.Rows
		err   error
		users []User
	)

	if rows, err = db.client.Query(getUserByAnyQuery, name, alias, email); err != nil {
		return nil, queryError(err, getUserByAnyQuery)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			println(fmt.Sprintf("error closing rows cause: %s", err.Error()))
		}
	}()

	for rows.Next() {
		var u User
		if err = rows.Scan(
			&u.ID,
			&u.UserName,
			&u.Alias,
			&u.Email,
			&u.Active,
			&u.DateCreated,
		); err != nil {
			return nil, scanError(err, getUserByAnyQuery)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, rowsError(err, getUserByAnyQuery)
	}

	return users, nil
}

func (db *relationalDB) CreateUser(request NewUserRequest) (User, error) {
	var (
		userID int64
		result sql.Result
		err    error

		dateCreated = time.Now()
	)

	result, err = db.client.Exec(insertUserQuery, request.UserName, request.Alias, request.Email, dateCreated)
	if err != nil {
		return User{}, execError(err, insertUserQuery)
	}

	if userID, err = result.LastInsertId(); err != nil {
		return User{}, lastInsertedError(err, insertUserQuery)
	}

	return User{
		ID:          userID,
		UserName:    request.UserName,
		Alias:       request.Alias,
		Email:       request.Email,
		DateCreated: dateCreated,
	}, nil
}
