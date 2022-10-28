package main

import (
	"github.com/gin-gonic/gin"
	"maria/src/api/user"
)

func main() {
	router := gin.Default()

	userController := user.NewController()

	userController.SetURLMapping(router)

	_ = router.Run("localhost:8080")
}
