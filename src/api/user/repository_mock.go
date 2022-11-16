package user

import (
	"errors"
	"fmt"

	"github.com/stretchr/testify/mock"
)

type dbMock struct {
	mock.Mock
}

func newDBMock() *dbMock {
	return &dbMock{}
}

func (m *dbMock) user(args mock.Arguments, index int) User {
	obj := args.Get(index)
	var s User
	var ok bool
	if s, ok = obj.(User); !ok {
		panic(fmt.Sprintf("assert: arguments: User(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
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
	return m.user(args, 0), args.Error(1)
}

func (m *dbMock) update(user User) (User, error) {
	args := m.Called(user)
	return m.user(args, 0), args.Error(1)
}

func (m *dbMock) insert(user User) (User, error) {
	args := m.Called(user)
	return m.user(args, 0), args.Error(1)
}
