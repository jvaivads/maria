package user

import (
	"errors"
	"testing"

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

func (s *UserServiceSuite) TestUserService() {
	var (
		userID = int64(10)
	)

	type test struct {
		name           string
		service        userService
		applyMockCalls func(us *userService) (func(t *testing.T), error)
		expectedError  error
		expectedUser   User
	}

	tests := []test{
		{
			name:           "user not found",
			service:        NewService(newDBMock()).(userService),
			applyMockCalls: setPersiterSelectByIDMock(User{}, nil, userID),
			expectedError:  userNotFoundByIDError,
			expectedUser:   User{},
		},
		{
			name:           "repository return error",
			service:        NewService(newDBMock()).(userService),
			applyMockCalls: setPersiterSelectByIDMock(User{}, errors.New("custom error"), userID),
			expectedError:  errors.New("custom error"),
			expectedUser:   User{},
		},
		{
			name:           "happy case",
			service:        NewService(newDBMock()).(userService),
			applyMockCalls: setPersiterSelectByIDMock(User{ID: userID}, nil, userID),
			expectedError:  nil,
			expectedUser:   User{ID: userID},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			if test.applyMockCalls != nil {
				if assertsCalls, err := test.applyMockCalls(&test.service); err != nil {
					assert.Fail(t, err.Error())
					return
				} else {
					defer assertsCalls(t)
				}
			}

			user, err := test.service.getByID(userID)

			assert.Equal(s.T(), test.expectedError, err)
			assert.Equal(s.T(), test.expectedUser, user)
		})
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
		r.On("SelectByID", userID).
			Return(userResponse, errorResponse).
			Once()
		return func(t *testing.T) {
			r.AssertExpectations(t)
		}, nil
	}
}
