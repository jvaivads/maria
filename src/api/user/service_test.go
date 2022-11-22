package user

import (
	"errors"
	"maria/src/api/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserServiceSuite struct {
	suite.Suite
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) BeforeTest(suiteName, testName string) {
}

func (s *UserServiceSuite) AfterTest(suiteName, testName string) {
}

func (s *UserServiceSuite) TestGetByID() {
	var (
		userID = int64(10)
	)

	type test struct {
		name           string
		applyMockCalls func(us *userService) (func(t *testing.T), error)
		expectedError  error
		expectedUser   User
	}

	tests := []test{
		{
			name:           "user not found",
			applyMockCalls: setPersiterSelectByIDMock(User{}, nil, userID),
			expectedError:  userNotFoundByIDError,
			expectedUser:   User{},
		},
		{
			name:           "repository return error",
			applyMockCalls: setPersiterSelectByIDMock(User{}, errors.New("custom error"), userID),
			expectedError:  errors.New("custom error"),
			expectedUser:   User{},
		},
		{
			name:           "happy case",
			applyMockCalls: setPersiterSelectByIDMock(User{ID: userID}, nil, userID),
			expectedError:  nil,
			expectedUser:   User{ID: userID},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			serv := NewService(newDBMock()).(userService)
			if test.applyMockCalls != nil {
				if assertsCalls, err := test.applyMockCalls(&serv); err != nil {
					assert.Fail(t, err.Error())
					return
				} else {
					defer assertsCalls(t)
				}
			}

			user, err := serv.getByID(userID)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedUser, user)
		})
	}
}

func (s *UserServiceSuite) TestCreateUser() {
	var (
		userID      = int64(10)
		customError = errors.New("custom error")
		userRequest = NewUserRequest{
			UserName: "name",
			Alias:    "alias",
			Email:    "email@email.com",
		}
	)

	type test struct {
		name          string
		mockCalls     mockPersisterApplier
		expectedError error
		expectedUser  User
	}

	tests := []test{
		{
			name:          "select by any return error",
			mockCalls:     mockPersisterApplier{setPersiterSelectByAnyMock(nil, customError, userRequest)},
			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name:          "select by any return a user with same name",
			mockCalls:     mockPersisterApplier{setPersiterSelectByAnyMock([]User{{UserName: userRequest.UserName}}, nil, userRequest)},
			expectedError: userWithSameValueErrorFunc("user_name"),
			expectedUser:  User{},
		},
		{
			name:          "select by any return a user with same alias",
			mockCalls:     mockPersisterApplier{setPersiterSelectByAnyMock([]User{{Alias: userRequest.Alias}}, nil, userRequest)},
			expectedError: userWithSameValueErrorFunc("alias"),
			expectedUser:  User{},
		},
		{
			name:          "select by any return a user with same email",
			mockCalls:     mockPersisterApplier{setPersiterSelectByAnyMock([]User{{Email: userRequest.Email}}, nil, userRequest)},
			expectedError: userWithSameValueErrorFunc("email"),
			expectedUser:  User{},
		},
		{
			name: "create user return error",
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock(nil, nil, userRequest),
				setPersiterCreateUserMock(User{}, customError, userRequest),
			},
			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "create user return error",
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock(nil, nil, userRequest),
				setPersiterCreateUserMock(userRequest.toUser(userID, time.Time{}, false), nil, userRequest),
			},
			expectedError: nil,
			expectedUser:  userRequest.toUser(userID, time.Time{}, false),
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			serv := NewService(newDBMock()).(userService)
			if assertsCalls, err := test.mockCalls.apply(&serv); err != nil {
				assert.Fail(t, err.Error())
				return
			} else {
				defer assertsCalls(t)
			}

			user, err := serv.createUser(userRequest)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedUser, user)
		})
	}
}

type mockPersisterApplier []func(us *userService) (func(t *testing.T), error)

func (appliers mockPersisterApplier) apply(us *userService) (func(t *testing.T), error) {
	var assertCalls []func(t *testing.T)
	for i := range appliers {
		if assertCall, err := appliers[i](us); err != nil {
			return func(t *testing.T) {}, err
		} else {
			assertCalls = append(assertCalls, assertCall)
		}
	}
	return func(t *testing.T) {
		for i := range assertCalls {
			assertCalls[i](t)
		}
	}, nil
}

func setPersiterSelectByIDMock(
	userResponse User,
	errorResponse error,
	userID int64,
) func(us *userService) (func(t *testing.T), error) {
	return func(us *userService) (func(t *testing.T), error) {
		r, ok := us.userRepository.(*dbMock)
		if !ok {
			return nil, errors.New("it could not cast to mock repository")
		}
		r.On(util.GetFunctionName(r.selectByID), userID).
			Return(userResponse, errorResponse).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
}

func setPersiterSelectByAnyMock(
	users []User,
	err error,
	userRequest NewUserRequest,
) func(us *userService) (func(t *testing.T), error) {
	return func(us *userService) (func(t *testing.T), error) {
		r, ok := us.userRepository.(*dbMock)
		if !ok {
			return nil, errors.New("it could not cast to mock repository")
		}
		r.On(
			util.GetFunctionName(r.selectByAny),
			userRequest.UserName,
			userRequest.Alias,
			userRequest.Email).
			Return(users, err).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
}

func setPersiterCreateUserMock(
	userResponse User,
	err error,
	userRequest NewUserRequest,
) func(us *userService) (func(t *testing.T), error) {
	return func(us *userService) (func(t *testing.T), error) {
		r, ok := us.userRepository.(*dbMock)
		if !ok {
			return nil, errors.New("it could not cast to mock repository")
		}
		r.On(
			util.GetFunctionName(r.createUser), userRequest).
			Return(userResponse, err).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
}
