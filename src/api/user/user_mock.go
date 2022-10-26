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

func (m *dbMock) user(args mock.Arguments, index int) user {
	obj := args.Get(index)
	var s user
	var ok bool
	if s, ok = obj.(user); !ok {
		panic(fmt.Sprintf("assert: arguments: user(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}

func (m *dbMock) onWithError(times int, method string, arguments ...interface{}) error {
	errorExpected := errors.New("custom error")
	m.On(method, arguments...).
		Return(user{}, errorExpected).
		Times(times)

	return errorExpected
}

func (m *dbMock) selectByID(userID int64) (user, error) {
	args := m.Called(userID)
	return m.user(args, 0), args.Error(1)
}

func (m *dbMock) update(user user) (user, error) {
	args := m.Called(user)
	return m.user(args, 0), args.Error(1)
}

func (m *dbMock) insert(user user) (user, error) {
	args := m.Called(user)
	return m.user(args, 0), args.Error(1)
}
