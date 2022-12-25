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

func (m *dbMock) selectByAny(name, alias, email string) ([]User, error) {
	args := m.Called(name, alias, email)
	return mockUsers(args, 0), args.Error(1)
}

func (m *dbMock) createUser(request NewUserRequest) (int64, error) {
	args := m.Called(request)
	return mockInt64(args, 0), args.Error(1)
}

func (m *dbMock) modifyUser(request ModifyUserRequest, user User) (bool, error) {
	args := m.Called(request, user)
	return args.Bool(0), args.Error(1)
}

func (m *dbMock) withTransaction(fn func(tx Transactioner) error) error {
	args := m.Called(fn)

	if err := args.Error(0); err != nil {
		return err
	}

	return fn(m)
}

func (m *dbMock) commit() error {
	args := m.Called()
	return args.Error(1)
}

func (m *dbMock) rollback() error {
	args := m.Called()
	return args.Error(1)
}
