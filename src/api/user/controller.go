package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	userIDMissedError = newBadRequestResponse("user_id param is missed")
)

type Controller struct {
	service       Service
	integerParser func(s string, base int, bitSize int) (i int64, err error)
}

func NewController(service Service) Controller {
	return Controller{
		service:       service,
		integerParser: strconv.ParseInt,
	}
}

func (c Controller) GetByID(ctx *gin.Context) {
	param := ctx.Param("user_id")
	if param == "" {
		ctx.JSON(http.StatusBadRequest, userIDMissedError)
		return
	}

	userID, err := c.integerParser(param, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newBadRequestResponse(err.Error()))
		return
	}

	user, err := c.service.getByID(userID)
	if err != nil {
		if errors.Is(err, userNotFoundError) {
			ctx.JSON(http.StatusNotFound, newNotFoundError("user_id", userID))
			return
		}
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c Controller) Post(ctx *gin.Context) {
	var (
		userRequest NewUserRequest
	)

	err := ctx.BindJSON(&userRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newBadRequestResponse(err.Error()))
		return
	}

	user, err := c.service.createUser(userRequest)
	if err != nil {
		if errors.Is(err, userWithSameValueError) {
			ctx.JSON(http.StatusBadRequest, newBadRequestResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c Controller) Put(ctx *gin.Context) {
	var (
		userID      int64
		userRequest ModifyUserRequest
		err         error
	)

	if err = ctx.BindJSON(&userRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, newBadRequestResponse(err.Error()))
		return
	}

	strUserID, ok1 := ctx.GetQuery("user_id")
	userName, ok2 := ctx.GetQuery("user_name")
	userAlias, ok3 := ctx.GetQuery("alias")

	if !ok1 && !ok2 && !ok3 {
		ctx.JSON(http.StatusBadRequest, newBadRequestResponse(
			"user_id, user_name nor email were not specified in query string"))
		return
	}

	if (ok1 && ok2) || (ok2 && ok3) || (ok3 && ok1) {
		ctx.JSON(http.StatusBadRequest, newBadRequestResponse(
			"specify only one parameter (user_id, user_name or email)"))
		return
	}

	if userID, err = strconv.ParseInt(strUserID, 10, 64); ok1 && err != nil {
		ctx.JSON(http.StatusBadRequest, newBadRequestResponse(
			"user_id must be a positive integer"))
		return
	}

	if userRequest.isEmpty() {
		ctx.JSON(http.StatusBadRequest, newBadRequestResponse(
			"request does not specify a change to be applied"))
		return
	}

	user := User{
		ID:       userID,
		UserName: userName,
		Alias:    userAlias,
	}

	if user, err = c.service.modifyUser(userRequest, user); err != nil {
		if errors.Is(err, userNotFoundError) {
			ctx.JSON(http.StatusBadRequest, newBadRequestResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, newInternalServerError(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c Controller) SetURLMapping(router *gin.Engine) {
	router.GET("/user/:user_id", c.GetByID)
	router.POST("/user", c.Post)
	router.PUT("/user/:user_id", c.Put)
}

func newBadRequestResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"message":     message,
		"status_code": http.StatusBadRequest,
	}
}

func newNotFoundError(by string, id any) map[string]interface{} {
	return map[string]interface{}{
		"message":     "element not found",
		"by":          by,
		"id":          id,
		"status_code": http.StatusNotFound,
	}
}

func newInternalServerError(cause error) map[string]interface{} {
	return map[string]interface{}{
		"message":     "internal server error",
		"cause":       cause,
		"status_code": http.StatusInternalServerError,
	}
}
