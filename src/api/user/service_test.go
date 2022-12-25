package user

import (
	"errors"
	"fmt"
	"maria/src/api/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			expectedError:  userNotFoundError,
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
			name: "select by any return error",
			mockCalls: mockPersisterApplier{setPersiterSelectByAnyMock(
				nil,
				customError,
				userRequest.UserName,
				userRequest.Alias,
				userRequest.Email)},
			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "select by any return a user with same name",
			mockCalls: mockPersisterApplier{setPersiterSelectByAnyMock(
				[]User{{UserName: userRequest.UserName}},
				nil,
				userRequest.UserName,
				userRequest.Alias,
				userRequest.Email)},
			expectedError: userWithSameValueErrorFunc("user_name"),
			expectedUser:  User{},
		},
		{
			name: "select by any return a user with same alias",
			mockCalls: mockPersisterApplier{setPersiterSelectByAnyMock(
				[]User{{Alias: userRequest.Alias}},
				nil,
				userRequest.UserName,
				userRequest.Alias,
				userRequest.Email)},
			expectedError: userWithSameValueErrorFunc("alias"),
			expectedUser:  User{},
		},
		{
			name: "select by any return a user with same email",
			mockCalls: mockPersisterApplier{setPersiterSelectByAnyMock(
				[]User{{Email: userRequest.Email}},
				nil,
				userRequest.UserName,
				userRequest.Alias,
				userRequest.Email)},
			expectedError: userWithSameValueErrorFunc("email"),
			expectedUser:  User{},
		},
		{
			name: "with transaction return error",
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock(
					nil,
					nil,
					userRequest.UserName,
					userRequest.Alias,
					userRequest.Email),
				setPersiterWithTransactionMock(customError),
			},
			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "create user return error",
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock(
					nil,
					nil,
					userRequest.UserName,
					userRequest.Alias,
					userRequest.Email),
				setPersiterWithTransactionMock(nil),
				setPersiterCreateUserMock(0, customError, userRequest),
			},
			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "select by user id  return error",
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock(
					nil,
					nil,
					userRequest.UserName,
					userRequest.Alias,
					userRequest.Email),
				setPersiterWithTransactionMock(nil),
				setPersiterCreateUserMock(userID, nil, userRequest),
				setPersiterSelectByIDMock(User{}, customError, userID),
			},
			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "create user is ok",
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock(
					nil,
					nil,
					userRequest.UserName,
					userRequest.Alias,
					userRequest.Email),
				setPersiterWithTransactionMock(nil),
				setPersiterCreateUserMock(userID, nil, userRequest),
				setPersiterSelectByIDMock(userRequest.toUser(userID, time.Time{}, false), nil, userID),
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

func (s *UserServiceSuite) TestModifyUser() {
	var (
		userID      = int64(10)
		active      = false
		customError = errors.New("custom error")
		userRequest = ModifyUserRequest{
			Active: &active,
		}
	)

	type test struct {
		name          string
		user          User
		mockCalls     mockPersisterApplier
		expectedError error
		expectedUser  User
	}

	tests := []test{
		{
			name: "select by any return error",
			user: User{UserName: "name"},
			mockCalls: mockPersisterApplier{setPersiterSelectByAnyMock(
				nil,
				customError,
				"name",
				"",
				"")},
			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "select by any return empty",
			user: User{UserName: "name"},
			mockCalls: mockPersisterApplier{setPersiterSelectByAnyMock(
				nil,
				nil,
				"name",
				"",
				"")},
			expectedError: userNotFoundError,
			expectedUser:  User{},
		},
		{
			name: "select by any return more than one",
			user: User{UserName: "name"},
			mockCalls: mockPersisterApplier{setPersiterSelectByAnyMock(
				make([]User, 2),
				nil,
				"name",
				"",
				"")},
			expectedError: fmt.Errorf("%w: there is more than one user", conflictError),
			expectedUser:  User{},
		},
		{
			name: "with transaction return error",
			user: User{UserName: "name"},
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock([]User{{ID: userID}}, nil, "name", "", ""),
				setPersiterWithTransactionMock(customError),
			},

			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "modify user return error",
			user: User{UserName: "name"},
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock([]User{{ID: userID}}, nil, "name", "", ""),
				setPersiterWithTransactionMock(nil),
				setPersiterModifyUserMock(false, customError, userRequest, User{ID: userID}),
			},

			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "select user by id return error",
			user: User{UserName: "name"},
			mockCalls: mockPersisterApplier{
				setPersiterSelectByAnyMock([]User{{ID: userID}}, nil, "name", "", ""),
				setPersiterWithTransactionMock(nil),
				setPersiterModifyUserMock(true, nil, userRequest, User{ID: userID}),
				setPersiterSelectByIDMock(User{}, customError, userID),
			},

			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "get by id return error",
			user: User{ID: userID},
			mockCalls: mockPersisterApplier{
				setPersiterSelectByIDMock(User{}, customError, userID),
			},

			expectedError: customError,
			expectedUser:  User{},
		},
		{
			name: "return ok",
			user: User{ID: userID},
			mockCalls: mockPersisterApplier{
				setPersiterSelectByIDMock(User{ID: userID}, nil, userID),
				setPersiterWithTransactionMock(nil),
				setPersiterModifyUserMock(true, nil, userRequest, User{ID: userID}),
				setPersiterSelectByIDMock(User{ID: userID}, nil, userID),
			},

			expectedError: nil,
			expectedUser:  User{ID: userID},
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

			user, err := serv.modifyUser(userRequest, test.user)

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

func setPersiterWithTransactionMock(
	errorResponse error,
) func(us *userService) (func(t *testing.T), error) {
	return func(us *userService) (func(t *testing.T), error) {
		r, ok := us.userRepository.(*dbMock)
		if !ok {
			return nil, errors.New("it could not cast to mock repository")
		}
		r.On(util.GetFunctionName(r.withTransaction), mock.Anything).
			Return(errorResponse).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
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
	userName, alias, email string,
) func(us *userService) (func(t *testing.T), error) {
	return func(us *userService) (func(t *testing.T), error) {
		r, ok := us.userRepository.(*dbMock)
		if !ok {
			return nil, errors.New("it could not cast to mock repository")
		}
		r.On(
			util.GetFunctionName(r.selectByAny),
			userName,
			alias,
			email).
			Return(users, err).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
}

func setPersiterCreateUserMock(
	userIDResponse int64,
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
			Return(userIDResponse, err).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
}

func setPersiterModifyUserMock(
	response bool,
	err error,
	request ModifyUserRequest,
	user User,
) func(us *userService) (func(t *testing.T), error) {
	return func(us *userService) (func(t *testing.T), error) {
		r, ok := us.userRepository.(*dbMock)
		if !ok {
			return nil, errors.New("it could not cast to mock repository")
		}
		r.On(
			util.GetFunctionName(r.modifyUser), request, user).
			Return(response, err).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
}
