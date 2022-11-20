package user

import (
	"database/sql"
	"errors"
	"maria/src/api/db"
	"testing"

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
				getUserMockRows([]User{{ID: userID}}),
				getUserByIDQuery,
				nil,
				userID),
			expectedError: nil,
			expectedUser:  User{ID: userID},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			client, mock, err := sqlmock.New()
			if err != nil {
				assert.Fail(t, err.Error())
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

			assert.Equal(s.T(), test.expectedError, err)
			assert.Equal(s.T(), test.expectedUser, user)
		})
	}
}

func getUserMockRows(users []User) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"ID"})

	for _, user := range users {
		rows.AddRow(user.ID)
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

func setClientQueryMock(
	row *sql.Rows,
	err error,
	query string,
	params ...any,
) func(c *db.ClientMock) (func(t *testing.T), error) {
	return func(c *db.ClientMock) (func(t *testing.T), error) {
		c.On("Query", query, params).
			Return(row, err).
			Once()
		return func(t *testing.T) {
			c.AssertExpectations(t)
		}, nil
	}
}
