package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"maria/src/api/db"
	"testing"
)

type UserServiceSuite struct {
	suite.Suite
	userDBMock  *db.UserDBMock
	userService User
	userID      int64
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) BeforeTest(suiteName, testName string) {
	s.userDBMock = db.NewUserDBMock()
	s.userService = NewUserService(s.userDBMock)
	s.userID = 10
}

func (s *UserServiceSuite) AfterTest(suiteName, testName string) {
	s.userDBMock.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestServiceReturnError() {
	errExpected := s.userDBMock.OnWithError(1, "SelectByID", s.userID)

	_, err := s.userService.GetByID(s.userID)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errExpected, err)
}
