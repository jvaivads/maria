package user

import (
	"database/sql"
	"errors"
	"fmt"
	"maria/src/api/db"
)

const (
	getUserByIDQuery    = `SELECT user_id, user_name, alias, email, active, date_created FROM user WHERE id = ?`
	getUserByAnyQuery   = `SELECT user_id, user_name, alias, email, active, date_created FROM user WHERE user_name = ? OR alias = ? OR email = ?`
	insertUserQuery     = `INSERT INTO user (user_name, alias, email, active) VALUES (?, ?, ?, false, NOW())`
	UpdateUserByIDQuery = `UPDATE user SET active = ? WHERE id = ?`
)

type Querier interface {
	selectByID(int64) (User, error)
	selectByAny(string, string, string) ([]User, error)
	createUser(NewUserRequest) (int64, error)
	modifyUser(ModifyUserRequest, User) (bool, error)
}

type Persister interface {
	Querier
	withTransaction(fn func(tx Transactioner) error) error
}

type Transactioner interface {
	Querier
	commit() error
	rollback() error
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

func (r *relationalDB) createUser(request NewUserRequest) (int64, error) {
	var (
		userID int64
		result sql.Result
		err    error
	)

	result, err = r.client.Exec(insertUserQuery, request.UserName, request.Alias, request.Email)
	if err != nil {
		return 0, db.ExecError(err, insertUserQuery)
	}

	if userID, err = result.LastInsertId(); err != nil {
		return 0, db.LastInsertedError(err, insertUserQuery)
	}

	return userID, nil
}

func (r *relationalDB) modifyUser(request ModifyUserRequest, user User) (bool, error) {
	if request.Active != nil {
		user.Active = *request.Active
	}
	result, err := r.client.Exec(UpdateUserByIDQuery, user.Active, user.ID)
	if err != nil {
		return false, db.ExecError(err, UpdateUserByIDQuery)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, db.RowsAffectedError(err, UpdateUserByIDQuery)
	}

	return rowsAffected == 1, nil
}

func (r *relationalDB) getTransactioner() (Transactioner, error) {
	client, ok := r.client.(*sql.DB)
	if !ok {
		return nil, errors.New("persister cannot generate transactional db")
	}

	tx, err := client.Begin()
	if err != nil {
		return nil, fmt.Errorf("persister cannot generate transactional due to: %w", err)
	}

	return &transactionalDB{relationalDB: relationalDB{client: tx}, tx: tx}, nil
}

func (r *relationalDB) withTransaction(fn func(tx Transactioner) error) error {
	tx, err := r.getTransactioner()
	if err != nil {
		return err
	}

	if err = fn(tx); err != nil {
		if err := tx.rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.commit()
}

type transactionalDB struct {
	relationalDB
	tx *sql.Tx
}

func (tx *transactionalDB) commit() error {
	if err := tx.tx.Commit(); err != nil {
		return db.CommitError(err)
	}
	return nil
}

func (tx *transactionalDB) rollback() error {
	if err := tx.tx.Rollback(); err != nil {
		return db.RollbackError(err)
	}
	return nil
}
