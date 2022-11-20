package user

import (
	"errors"
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

func (s *relationalDBSuite) TestRelationalDBSuite() {
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
			applyMockCalls: setClientQueryRowMock(
				getUserMockRows(nil),
				getUserByIDQuery,
				nil,
				userID),
			expectedError: nil,
			expectedUser:  User{},
		},
		{
			name: "scan error",
			applyMockCalls: setClientQueryRowMock(
				getUserMockRows([]User{{ID: userID}}),
				getUserByIDQuery,
				errors.New("custom error"),
				userID),
			expectedError: scanError(errors.New("custom error"), getUserByIDQuery),
			expectedUser:  User{},
		},
		{
			name: "happy case",
			applyMockCalls: setClientQueryRowMock(
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

			user, err := rDB.SelectByID(userID)

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

func setClientQueryRowMock(
	rows *sqlmock.Rows,
	query string,
	scanError error,
	params ...any,
) func(m sqlmock.Sqlmock) func() error {
	return func(m sqlmock.Sqlmock) func() error {
		m.ExpectQuery(query).WithArgs(params[0]).WillReturnRows(rows)
		if scanError != nil {
			rows.RowError(0, scanError)
		}
		return func() error {
			return m.ExpectationsWereMet()
		}
	}
}
