package db

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"maria/src/api/domain"
)

type UserDBMock struct {
	mock.Mock
}

func (m *UserDBMock) user(args mock.Arguments, index int) domain.User {
	obj := args.Get(index)
	var s domain.User
	var ok bool
	if s, ok = obj.(domain.User); !ok {
		panic(fmt.Sprintf("assert: arguments: user(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}

func (m *UserDBMock) SelectByID(userID int64) (domain.User, error) {
	args := m.Called(userID)
	return m.user(args, 0), args.Error(1)
}

func (m *UserDBMock) Update(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return m.user(args, 0), args.Error(1)
}

func (m *UserDBMock) Insert(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return m.user(args, 0), args.Error(1)
}
