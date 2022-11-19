package user

import (
	"errors"

	"github.com/stretchr/testify/mock"
)

type dbMock struct {
	mock.Mock
}

func newDBMock() *dbMock {
	return &dbMock{}
}

func (m *dbMock) onWithError(times int, method string, arguments ...interface{}) error {
	errorExpected := errors.New("custom error")
	m.On(method, arguments...).
		Return(User{}, errorExpected).
		Times(times)

	return errorExpected
}

func (m *dbMock) SelectByID(userID int64) (User, error) {
	args := m.Called(userID)
	return mockUser(args, 0), args.Error(1)
}

func (m *dbMock) update(user User) (User, error) {
	args := m.Called(user)
	return mockUser(args, 0), args.Error(1)
}

func (m *dbMock) insert(user User) (User, error) {
	args := m.Called(user)
	return mockUser(args, 0), args.Error(1)
}
