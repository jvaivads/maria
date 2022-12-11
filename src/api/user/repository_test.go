package user

import (
	"errors"
	"maria/src/api/db"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type relationalDBSuite struct {
	suite.Suite
}

func TestRelationalDBSuite(t *testing.T) {
	suite.Run(t, new(relationalDBSuite))
}

func (s *relationalDBSuite) BeforeTest(suiteName, testName string) {
}

func (s *relationalDBSuite) AfterTest(suiteName, testName string) {
}

func (s *relationalDBSuite) TestSelectByID() {
	var (
		userID = int64(10)
		user   = User{
			ID:          userID,
			UserName:    "user",
			Alias:       "alias",
			Email:       "user@email.com",
			DateCreated: time.Now(),
			Active:      true,
		}
	)

	type test struct {
		name           string
		applyMockCalls func(m sqlmock.Sqlmock) func() error
		expectedError  error
		expectedUser   User
	}

	tests := []test{
		{
			name: "empty result",
			applyMockCalls: db.SetClientQueryRowMock(
				getUserMockRows(nil),
				getUserByIDQuery,
				nil,
				userID),
			expectedError: nil,
			expectedUser:  User{},
		},
		{
			name: "scan error",
			applyMockCalls: db.SetClientQueryRowMock(
				getUserMockRows([]User{{ID: userID}}),
				getUserByIDQuery,
				errors.New("custom error"),
				userID),
			expectedError: db.ScanError(errors.New("custom error"), getUserByIDQuery),
			expectedUser:  User{},
		},
		{
			name: "happy case",
			applyMockCalls: db.SetClientQueryRowMock(
				getUserMockRows([]User{user}),
				getUserByIDQuery,
				nil,
				userID),
			expectedError: nil,
			expectedUser:  user,
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if test.applyMockCalls != nil {
				assertsCalls := test.applyMockCalls(mock)
				defer func() {
					if err = assertsCalls(); err != nil {
						assert.Fail(t, err.Error())
					}
				}()
			}

			rDB := NewRelationalDB(client)

			user, err := rDB.selectByID(userID)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedUser, user)
		})
	}
}

func (s *relationalDBSuite) TestSelectByAny() {
	var (
		userID      = int64(10)
		userName    = "name"
		alias       = "alias"
		email       = "email@email.com"
		customError = errors.New("custom error")
		user        = User{
			ID:          userID,
			UserName:    "user",
			Alias:       "alias",
			Email:       "user@email.com",
			DateCreated: time.Now(),
			Active:      true,
		}
	)

	type test struct {
		name           string
		applyMockCalls func(m sqlmock.Sqlmock) func() error
		expectedError  error
		expectedUsers  []User
	}

	tests := []test{
		{
			name: "query error",
			applyMockCalls: db.SetClientQueryMock(
				getUserMockRows(nil),
				getUserByAnyQuery,
				customError,
				nil,
				userName, alias, email),
			expectedError: db.QueryError(customError, getUserByAnyQuery),
			expectedUsers: nil,
		},
		{
			name: "rows error",
			applyMockCalls: db.SetClientQueryMock(
				getUserMockRows([]User{user}),
				getUserByAnyQuery,
				nil,
				customError,
				userName, alias, email),
			expectedError: db.RowsError(customError, getUserByAnyQuery),
			expectedUsers: nil,
		},
		{
			name: "rows error",
			applyMockCalls: db.SetClientQueryMock(
				getUserMockRows([]User{user}),
				getUserByAnyQuery,
				nil,
				nil,
				userName, alias, email),
			expectedError: nil,
			expectedUsers: []User{user},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if test.applyMockCalls != nil {
				assertsCalls := test.applyMockCalls(mock)
				defer func() {
					if err = assertsCalls(); err != nil {
						assert.Fail(t, err.Error())
					}
				}()
			}

			rDB := NewRelationalDB(client)

			user, err := rDB.selectByAny(userName, alias, email)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedUsers, user)
		})
	}
}

func (s *relationalDBSuite) TestCreateUser() {
	var (
		userID      = int64(10)
		dateCreated = time.Now()
		userRequest = NewUserRequest{
			UserName: "name",
			Alias:    "alias",
			Email:    "email@email.com",
		}
		customError = errors.New("custom error")
	)

	type test struct {
		name          string
		mockCalls     mockDBApplier
		expectedError error
		expectedUser  User
	}

	tests := []test{
		{
			name: "query error",
			mockCalls: mockDBApplier{db.SetClientExecMock(
				nil,
				insertUserQuery,
				customError,
				userRequest.UserName, userRequest.Alias, userRequest.Email),
			},
			expectedError: db.ExecError(customError, insertUserQuery),
			expectedUser:  User{},
		},
		{
			name: "last inserted error",
			mockCalls: mockDBApplier{db.SetClientExecMock(
				sqlmock.NewErrorResult(customError),
				insertUserQuery,
				nil,
				userRequest.UserName, userRequest.Alias, userRequest.Email),
			},
			expectedError: db.LastInsertedError(customError, insertUserQuery),
			expectedUser:  User{},
		},
		{
			name: "happy case",
			mockCalls: mockDBApplier{
				db.SetClientExecMock(
					sqlmock.NewResult(10, 1),
					insertUserQuery,
					nil,
					userRequest.UserName, userRequest.Alias, userRequest.Email),
				db.SetClientQueryMock(
					getUserMockRows([]User{userRequest.toUser(userID, dateCreated, false)}),
					getUserByIDQuery,
					nil,
					nil,
					userID),
			},
			expectedError: nil,
			expectedUser:  userRequest.toUser(userID, dateCreated, false),
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assertsCalls := test.mockCalls.apply(mock)
			defer func() {
				if err = assertsCalls(); err != nil {
					assert.Fail(t, err.Error())
				}
			}()

			rDB := NewRelationalDB(client)

			user, err := rDB.createUser(userRequest)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedUser, user)
		})
	}
}

func (s *relationalDBSuite) TestModifyUser() {
	var (
		user        = User{ID: 10}
		active      = true
		userRequest = ModifyUserRequest{Active: &active}
		customError = errors.New("custom error")
	)

	type test struct {
		name          string
		mockCalls     mockDBApplier
		expectedError error
		expectedUser  User
	}

	tests := []test{
		{
			name: "query error",
			mockCalls: mockDBApplier{db.SetClientExecMock(
				nil,
				UpdateUserByIDQuery,
				customError,
				active, user.ID),
			},
			expectedError: db.ExecError(customError, UpdateUserByIDQuery),
			expectedUser:  User{},
		},
		{
			name: "rows affected error",
			mockCalls: mockDBApplier{db.SetClientExecMock(
				sqlmock.NewErrorResult(customError),
				UpdateUserByIDQuery,
				nil,
				active, user.ID),
			},
			expectedError: db.RowsAffectedError(customError, UpdateUserByIDQuery),
			expectedUser:  User{},
		},
		{
			name: "select by id return error",
			mockCalls: mockDBApplier{
				db.SetClientExecMock(
					sqlmock.NewResult(0, 1),
					UpdateUserByIDQuery,
					nil,
					active, user.ID),
				db.SetClientQueryRowMock(
					getUserMockRows([]User{{ID: user.ID}}),
					getUserByIDQuery,
					errors.New("custom error"),
					user.ID),
			},
			expectedError: db.ScanError(customError, getUserByIDQuery),
			expectedUser:  User{},
		},
		{
			name: "happy case",
			mockCalls: mockDBApplier{
				db.SetClientExecMock(
					sqlmock.NewResult(0, 1),
					UpdateUserByIDQuery,
					nil,
					active, user.ID),
				db.SetClientQueryRowMock(
					getUserMockRows([]User{{ID: user.ID}}),
					getUserByIDQuery,
					nil,
					user.ID),
			},
			expectedError: nil,
			expectedUser:  User{ID: user.ID},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			assertsCalls := test.mockCalls.apply(mock)
			defer func() {
				if err = assertsCalls(); err != nil {
					assert.Fail(t, err.Error())
				}
			}()

			rDB := NewRelationalDB(client)

			user, err := rDB.modifyUser(userRequest, user)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedUser, user)
		})
	}
}

func getUserMockRows(users []User) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"user_id", "user_name", "alias", "email", "active", "date_created"})

	for _, user := range users {
		rows.AddRow(user.ID, user.UserName, user.Alias, user.Email, user.Active, user.DateCreated)
	}
	return rows
}

type mockDBApplier []func(m sqlmock.Sqlmock) func() error

func (appliers mockDBApplier) apply(m sqlmock.Sqlmock) func() error {
	var assertCalls []func() error
	for i := range appliers {
		assertCall := appliers[i](m)
		assertCalls = append(assertCalls, assertCall)
	}
	return func() error {
		for i := range assertCalls {
			if err := assertCalls[i](); err != nil {
				return err
			}
		}
		return nil
	}
}
