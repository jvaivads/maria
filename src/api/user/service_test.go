package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserServiceSuite struct {
	suite.Suite
	dbMock  *dbMock
	service service
	userID  int64
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) BeforeTest(suiteName, testName string) {
	s.dbMock = newDBMock()
	s.service = newService(s.dbMock)
	s.userID = 10
}

func (s *UserServiceSuite) AfterTest(suiteName, testName string) {
	s.dbMock.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestServiceReturnError() {
	errExpected := s.dbMock.onWithError(1, "selectByID", s.userID)

	_, err := s.service.getByID(s.userID)

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), errExpected, err)
}
