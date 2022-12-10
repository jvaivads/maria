package user

import (
	"database/sql"
	"fmt"
	"maria/src/api/db"
)

const (
	getUserByIDQuery    = `SELECT user_id, user_name, alias, email, active, date_created FROM user WHERE id = ?`
	getUserByAnyQuery   = `SELECT user_id, user_name, alias, email, active, date_created FROM user WHERE user_name = ? OR alias = ? OR email = ?`
	insertUserQuery     = `INSERT INTO user (user_name, alias, email, active) VALUES (?, ?, ?, false, NOW())`
	UpdateUserByIDQuery = `UPDATE user SET active = ? WHERE id = ?`
)

type Persister interface {
	selectByID(int64) (User, error)
	selectByAny(string, string, string) ([]User, error)
	createUser(NewUserRequest) (User, error)
	modifyUser(ModifyUserRequest, User) (User, error)
}

func NewRelationalDB(client db.Client) Persister {
	return &relationalDB{
		client: client,
	}
}

type relationalDB struct {
	client db.Client
}

func (r *relationalDB) selectByID(userID int64) (User, error) {
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

func (r *relationalDB) selectByAny(name, alias, email string) ([]User, error) {
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

func (r *relationalDB) createUser(request NewUserRequest) (User, error) {
	var (
		userID int64
		result sql.Result
		err    error
	)

	result, err = r.client.Exec(insertUserQuery, request.UserName, request.Alias, request.Email)
	if err != nil {
		return User{}, db.ExecError(err, insertUserQuery)
	}

	if userID, err = result.LastInsertId(); err != nil {
		return User{}, db.LastInsertedError(err, insertUserQuery)
	}

	return r.selectByID(userID)
}

func (r *relationalDB) modifyUser(request ModifyUserRequest, user User) (User, error) {
	result, err := r.client.Exec(UpdateUserByIDQuery, request.Active, user.ID)
	if err != nil {
		return User{}, db.ExecError(err, UpdateUserByIDQuery)
	}

	if _, err := result.RowsAffected(); err != nil {
		return User{}, db.RowsAffectedError(err, UpdateUserByIDQuery)
	}

	return r.selectByID(user.ID)
}
