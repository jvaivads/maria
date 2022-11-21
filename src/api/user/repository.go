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

func (r *relationalDB) SelectByID(userID int64) (User, error) {
	var (
		u   User
		row *sql.Row
		err error
	)

	row = r.client.QueryRow(getUserByIDQuery, userID)

	if err = row.Scan(
		&u.ID,
		&u.UserName,
		&u.Alias,
		&u.Email,
		&u.Active,
		&u.DateCreated,
	); err != nil {
		return u, db.ScanError(err, getUserByIDQuery)
	}

	return u, nil
}

func (r *relationalDB) SelectByAny(name, alias, email string) ([]User, error) {
	var (
		rows  *sql.Rows
		err   error
		users []User
	)

	if rows, err = r.client.Query(getUserByAnyQuery, name, alias, email); err != nil {
		return nil, db.QueryError(err, getUserByAnyQuery)
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
			return nil, db.ScanError(err, getUserByAnyQuery)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, db.RowsError(err, getUserByAnyQuery)
	}

	return users, nil
}

func (r *relationalDB) CreateUser(request NewUserRequest) (User, error) {
	var (
		userID int64
		result sql.Result
		err    error

		dateCreated = time.Now()
	)

	result, err = r.client.Exec(insertUserQuery, request.UserName, request.Alias, request.Email, dateCreated)
	if err != nil {
		return User{}, db.ExecError(err, insertUserQuery)
	}

	if userID, err = result.LastInsertId(); err != nil {
		return User{}, db.LastInsertedError(err, insertUserQuery)
	}

	return User{
		ID:          userID,
		UserName:    request.UserName,
		Alias:       request.Alias,
		Email:       request.Email,
		DateCreated: dateCreated,
	}, nil
}
