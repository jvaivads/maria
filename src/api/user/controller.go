package user

import (
	"fmt"
	"html"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID int64
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{service: service}
}

func (c Controller) GetByID(ctx *gin.Context) {
	_, _ = fmt.Fprintf(ctx.Writer, "Hello, %q", html.EscapeString(ctx.Request.URL.Path))
}

func (c Controller) SetURLMapping(router *gin.Engine) {
	router.GET("/users/:user_id", c.GetByID)
}
