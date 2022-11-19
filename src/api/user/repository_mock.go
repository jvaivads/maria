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
