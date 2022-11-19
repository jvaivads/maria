package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ControllerSuite struct {
	suite.Suite
}

func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(ControllerSuite))
}

func (s *ControllerSuite) BeforeTest(suiteName, testName string) {

}

func (s *ControllerSuite) AfterTest(suiteName, testName string) {
}

func (c *ControllerSuite) TestGetUserByID() {
	var (
		userID      = int64(10)
		customError = errors.New("custom error")
	)
	type test struct {
		name           string
		param          string
		controller     Controller
		applyMockCalls func(controller *Controller) error
		expectedCode   int
		expectedBody   map[string]interface{}
	}

	tests := []test{
		{
			name:         "user_id param missed",
			param:        "",
			controller:   Controller{},
			expectedCode: http.StatusBadRequest,
			expectedBody: newBadRequestResponse("user_id param is missed"),
		},
		{
			name:  "it cannot parse user_id param",
			param: "word",
			controller: Controller{
				integerParser: func(s string, base int, bitSize int) (i int64, err error) {
					return 0, customError
				},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: newBadRequestResponse(customError.Error()),
		},
		{
			name:       "user not found",
			param:      "10",
			controller: NewController(newServiceMock()),
			applyMockCalls: func(controller *Controller) error {
				s, ok := controller.service.(*serviceMock)
				if !ok {
					return errors.New("it could not cast to mock service")
				}
				s.On("getByID", userID).
					Return(User{}, userNotFoundByIDError).
					Once()
				return nil
			},
			expectedCode: http.StatusNotFound,
			expectedBody: newNotFoundError("user_id", userID),
		},
		{
			name:       "service return internal server error",
			param:      "10",
			controller: NewController(newServiceMock()),
			applyMockCalls: func(controller *Controller) error {
				s, ok := controller.service.(*serviceMock)
				if !ok {
					return errors.New("it could not cast to mock service")
				}
				s.On("getByID", userID).
					Return(User{}, customError).
					Once()
				return nil
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: newInternalServerError(customError),
		},
		{
			name:       "happy case",
			param:      "10",
			controller: NewController(newServiceMock()),
			applyMockCalls: func(controller *Controller) error {
				s, ok := controller.service.(*serviceMock)
				if !ok {
					return errors.New("it could not cast to mock service")
				}
				s.On("getByID", int64(10)).
					Return(User{ID: userID}, nil).
					Once()
				return nil
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{"user_id": userID},
		},
	}

	for _, test := range tests {
		c.T().Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(r)
			ctx.AddParam("user_id", test.param)

			if test.applyMockCalls != nil {
				if err := test.applyMockCalls(&test.controller); err != nil {
					assert.Fail(t, err.Error())
					return
				}
			}

			test.controller.GetUserByID(ctx)

			bytes, err := json.Marshal(test.expectedBody)
			if !assert.Nil(t, err) {
				return
			}

			assert.Equal(t, test.expectedCode, r.Code)
			assert.Equal(t, string(bytes), r.Body.String())
		})
	}
}
