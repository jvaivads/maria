package user

import (
	"fmt"

	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func newServiceMock() *serviceMock {
	return &serviceMock{}
}

func mockUser(args mock.Arguments, index int) User {
	obj := args.Get(index)
	var s User
	var ok bool
	if s, ok = obj.(User); !ok {
		panic(fmt.Sprintf("assert: arguments: User(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}

func mockUsers(args mock.Arguments, index int) []User {
	obj := args.Get(index)
	var s []User
	var ok bool
	if s, ok = obj.([]User); !ok {
		panic(fmt.Sprintf("assert: arguments: User(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}

func (m *serviceMock) getByID(userID int64) (User, error) {
	args := m.Called(userID)
	return mockUser(args, 0), args.Error(1)
}

func (m *serviceMock) createUser(user NewUserRequest) (User, error) {
	args := m.Called(user)
	return mockUser(args, 0), args.Error(1)
}
