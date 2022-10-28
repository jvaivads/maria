package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
)

type user struct {
	ID int64
}

type Controller struct {
	service service
}

func NewController() Controller {
	return Controller{service: newService(newRelationalDB())}
}

func (c Controller) GetByID(ctx *gin.Context) {
	_, _ = fmt.Fprintf(ctx.Writer, "Hello, %q", html.EscapeString(ctx.Request.URL.Path))
}

func (c Controller) SetURLMapping(router *gin.Engine) {
	router.GET("/users/:user_id", c.GetByID)
}
