package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maria/src/api/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func (c *ControllerSuite) BeforeTest(suiteName, testName string) {

}

func (c *ControllerSuite) AfterTest(suiteName, testName string) {
}

func (c *ControllerSuite) TestGetUserByID() {
	var (
		userID      = int64(10)
		customError = errors.New("custom error")
		user        = User{
			ID:          userID,
			UserName:    "user",
			Alias:       "alias",
			Email:       "user@email.com",
			DateCreated: time.Now(),
			Active:      true,
		}
	)

	type test struct {
		name           string
		param          string
		controller     Controller
		applyMockCalls func(controller *Controller) (func(t *testing.T), error)
		expectedCode   int
		expectedBody   string
	}

	tests := []test{
		{
			name:         "user_id param missed",
			param:        "",
			controller:   Controller{},
			expectedCode: http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse("user_id param is missed")),
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
			expectedBody: util.RenderToJSON(newBadRequestResponse(customError.Error())),
		},
		{
			name:           "user not found",
			param:          "10",
			controller:     NewController(newServiceMock()),
			applyMockCalls: setServiceGetByIDMock(User{}, userNotFoundError, userID),
			expectedCode:   http.StatusNotFound,
			expectedBody:   util.RenderToJSON(newNotFoundError("user_id", userID)),
		},
		{
			name:           "service return internal server error",
			param:          "10",
			controller:     NewController(newServiceMock()),
			applyMockCalls: setServiceGetByIDMock(User{}, customError, userID),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   util.RenderToJSON(newInternalServerError(customError)),
		},
		{
			name:           "happy case",
			param:          "10",
			controller:     NewController(newServiceMock()),
			applyMockCalls: setServiceGetByIDMock(user, nil, userID),
			expectedCode:   http.StatusOK,
			expectedBody:   util.RenderToJSON(user),
		},
	}

	for _, test := range tests {
		c.T().Run(test.name, func(t *testing.T) {
			params := map[string]string{
				"user_id": test.param,
			}
			ctx, r, err := util.GetTestContext(params, "", nil)
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}
			ctx.AddParam("user_id", test.param)

			if test.applyMockCalls != nil {
				if assertsCalls, err := test.applyMockCalls(&test.controller); err != nil {
					assert.Fail(t, err.Error())
					return
				} else {
					defer assertsCalls(t)
				}
			}

			test.controller.GetByID(ctx)

			assert.Equal(t, test.expectedCode, r.Code)
			assert.Equal(t, test.expectedBody, r.Body.String())
		})
	}
}

func (c *ControllerSuite) TestPost() {
	const (
		bindMSGError = "" +
			"Key: 'NewUserRequest.UserName' Error:Field validation for 'UserName' failed on the 'required' tag\n" +
			"Key: 'NewUserRequest.Alias' Error:Field validation for 'Alias' failed on the 'required' tag\n" +
			"Key: 'NewUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
	)
	var (
		userID         = int64(10)
		dateCreated    = time.Now()
		customError    = errors.New("custom error")
		sameValueError = fmt.Errorf("custom error cause: %w", userWithSameValueError)

		userRequest = NewUserRequest{
			UserName: "name",
			Alias:    "alias",
			Email:    "email@email.com",
		}
	)

	type test struct {
		name           string
		body           NewUserRequest
		controller     Controller
		applyMockCalls func(controller *Controller) (func(t *testing.T), error)
		expectedCode   int
		expectedBody   string
	}

	tests := []test{
		{
			name:         "body fields missed",
			body:         NewUserRequest{},
			controller:   Controller{},
			expectedCode: http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse(bindMSGError)),
		},
		{
			name:       "service return user with same value error",
			body:       userRequest,
			controller: NewController(newServiceMock()),
			applyMockCalls: setServicePostMock(
				User{},
				sameValueError,
				userRequest,
			),
			expectedCode: http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse(sameValueError.Error())),
		},
		{
			name:       "service return internal error",
			body:       userRequest,
			controller: NewController(newServiceMock()),
			applyMockCalls: setServicePostMock(
				User{},
				customError,
				userRequest,
			),
			expectedCode: http.StatusInternalServerError,
			expectedBody: util.RenderToJSON(newInternalServerError(customError)),
		},
		{
			name:       "happy case",
			body:       userRequest,
			controller: NewController(newServiceMock()),
			applyMockCalls: setServicePostMock(
				userRequest.toUser(userID, dateCreated, false),
				nil,
				userRequest,
			),
			expectedCode: http.StatusOK,
			expectedBody: util.RenderToJSON(userRequest.toUser(userID, dateCreated, false)),
		},
	}

	for _, test := range tests {
		c.T().Run(test.name, func(t *testing.T) {
			ctx, r, err := util.GetTestContext(nil, "", test.body)
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if test.applyMockCalls != nil {
				if assertsCalls, err := test.applyMockCalls(&test.controller); err != nil {
					assert.Fail(t, err.Error())
					return
				} else {
					defer assertsCalls(t)
				}
			}

			test.controller.Post(ctx)

			assert.Equal(t, test.expectedCode, r.Code)
			assert.Equal(t, test.expectedBody, r.Body.String())
		})
	}
}

func (c *ControllerSuite) TestPut() {
	var (
		active          = true
		requestToActive = ModifyUserRequest{Active: &active}
		userID          = int64(10)
		customError     = errors.New("custom error")
	)

	type test struct {
		name           string
		queryString    string
		body           ModifyUserRequest
		controller     Controller
		applyMockCalls func(controller *Controller) (func(t *testing.T), error)
		expectedCode   int
		expectedBody   string
	}

	tests := []test{
		{
			name:           "not enough data",
			body:           ModifyUserRequest{},
			controller:     NewController(newServiceMock()),
			applyMockCalls: nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse(
				"user_id, user_name nor email were not specified in query string")),
		},
		{
			name:           "more than one query parameter",
			body:           ModifyUserRequest{},
			queryString:    "user_id=value&user_name=value",
			controller:     NewController(newServiceMock()),
			applyMockCalls: nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse(
				"specify only one parameter (user_id, user_name or email)")),
		},
		{
			name:           "user_id is not a positive integer",
			body:           ModifyUserRequest{},
			queryString:    "user_id=value",
			controller:     NewController(newServiceMock()),
			applyMockCalls: nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse(
				"user_id must be a positive integer")),
		},
		{
			name:           "request is empty",
			body:           ModifyUserRequest{},
			queryString:    "user_id=10",
			controller:     NewController(newServiceMock()),
			applyMockCalls: nil,
			expectedCode:   http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse(
				"request does not specify a change to be applied")),
		},
		{
			name:        "user not found",
			body:        requestToActive,
			queryString: "user_id=10",
			controller:  NewController(newServiceMock()),
			applyMockCalls: setServicePutMock(
				User{},
				fmt.Errorf("error: %w", userNotFoundError),
				requestToActive,
				User{ID: userID}),
			expectedCode: http.StatusBadRequest,
			expectedBody: util.RenderToJSON(newBadRequestResponse(
				fmt.Errorf("error: %w", userNotFoundError).Error())),
		},
		{
			name:           "internal error",
			body:           requestToActive,
			queryString:    "user_name=name",
			controller:     NewController(newServiceMock()),
			applyMockCalls: setServicePutMock(User{}, customError, requestToActive, User{UserName: "name"}),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   util.RenderToJSON(newInternalServerError(customError)),
		},
		{
			name:           "internal error",
			body:           requestToActive,
			queryString:    "alias=alias",
			controller:     NewController(newServiceMock()),
			applyMockCalls: setServicePutMock(User{Alias: "alias"}, nil, requestToActive, User{Alias: "alias"}),
			expectedCode:   http.StatusOK,
			expectedBody:   util.RenderToJSON(User{Alias: "alias"}),
		},
	}

	for _, test := range tests {
		c.T().Run(test.name, func(t *testing.T) {
			ctx, r, err := util.GetTestContext(nil, test.queryString, test.body)
			if err != nil {
				assert.Fail(t, err.Error())
				return
			}

			if test.applyMockCalls != nil {
				if assertsCalls, err := test.applyMockCalls(&test.controller); err != nil {
					assert.Fail(t, err.Error())
					return
				} else {
					defer assertsCalls(t)
				}
			}

			test.controller.Put(ctx)

			assert.Equal(t, test.expectedCode, r.Code)
			assert.Equal(t, test.expectedBody, r.Body.String())
		})
	}
}

func (c *ControllerSuite) TestSetURLMapping() {
	var (
		userID      = int64(10)
		userRequest = NewUserRequest{
			UserName: "name",
			Alias:    "alias",
			Email:    "email@email.com",
		}
	)

	type test struct {
		name           string
		path           string
		method         string
		body           any
		controller     Controller
		applyMockCalls func(controller *Controller) (func(t *testing.T), error)
	}

	tests := []test{
		{
			name:           "get user by id",
			path:           "/user/10",
			method:         http.MethodGet,
			controller:     NewController(newServiceMock()),
			applyMockCalls: setServiceGetByIDMock(User{ID: userID}, nil, userID),
		},
		{
			name:       "post user",
			path:       "/user",
			method:     http.MethodPost,
			body:       userRequest,
			controller: NewController(newServiceMock()),
			applyMockCalls: setServicePostMock(
				userRequest.toUser(userID, time.Time{}, false),
				nil,
				userRequest,
			),
		},
	}

	for _, test := range tests {
		c.T().Run(test.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			test.controller.SetURLMapping(router)

			if test.applyMockCalls != nil {
				if assertsCalls, err := test.applyMockCalls(&test.controller); err != nil {
					assert.Fail(t, err.Error())
					return
				} else {
					defer assertsCalls(t)
				}
			}

			var b bytes.Buffer
			if test.body != nil {
				if err := json.NewEncoder(&b).Encode(test.body); err != nil {
					assert.Fail(t, err.Error())
					return
				}
			}

			req := httptest.NewRequest(test.method, test.path, &b)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(c.T(), http.StatusOK, w.Code)
		})
	}
}

func setServiceGetByIDMock(
	userResponse User,
	errorResponse error,
	userID int64,
) func(*Controller) (func(t *testing.T), error) {
	return func(c *Controller) (func(t *testing.T), error) {
		s, ok := c.service.(*serviceMock)
		if !ok {
			return nil, errors.New("it could not cast to mock service")
		}
		s.On(util.GetFunctionName(s.getByID), userID).
			Return(userResponse, errorResponse).
			Once()

		return func(t *testing.T) {
			s.AssertExpectations(t)
		}, nil
	}
}

func setServicePostMock(
	userResponse User,
	errorResponse error,
	userRequest NewUserRequest,
) func(*Controller) (func(t *testing.T), error) {
	return func(c *Controller) (func(t *testing.T), error) {
		s, ok := c.service.(*serviceMock)
		if !ok {
			return nil, errors.New("it could not cast to mock service")
		}
		s.On(util.GetFunctionName(s.createUser), userRequest).
			Return(userResponse, errorResponse).
			Once()

		return func(t *testing.T) {
			s.AssertExpectations(t)
		}, nil
	}
}

func setServicePutMock(
	userResponse User,
	errorResponse error,
	request ModifyUserRequest,
	userRequest User,
) func(*Controller) (func(t *testing.T), error) {
	return func(c *Controller) (func(t *testing.T), error) {
		s, ok := c.service.(*serviceMock)
		if !ok {
			return nil, errors.New("it could not cast to mock service")
		}
		s.On(util.GetFunctionName(s.modifyUser), request, userRequest).
			Return(userResponse, errorResponse).
			Once()

		return func(t *testing.T) {
			s.AssertExpectations(t)
		}, nil
	}
}
