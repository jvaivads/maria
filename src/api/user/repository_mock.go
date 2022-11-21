package user

import (
	"github.com/stretchr/testify/mock"
)

type dbMock struct {
	mock.Mock
}

func newDBMock() *dbMock {
	return &dbMock{}
}

func (m *dbMock) selectByID(userID int64) (User, error) {
	args := m.Called(userID)
	return mockUser(args, 0), args.Error(1)
}

func (m *dbMock) SelectByAny(name, alias, email string) ([]User, error) {
	args := m.Called(name, alias, email)
	return mockUsers(args, 0), args.Error(1)
}

func (m *dbMock) CreateUser(request NewUserRequest) (User, error) {
	args := m.Called(request)
	return mockUser(args, 0), args.Error(1)
}
