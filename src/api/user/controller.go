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
		if errors.Is(err, userNotFoundByIDError) {
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
	}

	user, err := c.service.createUser(userRequest)
	if err != nil {
		if errors.Is(err, userWithSameValueError) {
			ctx.JSON(http.StatusBadRequest, err)
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
