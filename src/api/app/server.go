package main

import (
	"maria/src/api/user"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	userController := user.NewController()

	userController.SetURLMapping(router)

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
