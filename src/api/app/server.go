package main

import (
	"maria/src/api/db"
	"maria/src/api/user"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	userController := user.NewController(
		user.NewService(
			user.NewRelationalDB(
				db.NewSQLClient(
					getSQLClientConfig(),
				),
			),
		),
	)

	userController.SetURLMapping(router)

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}

func getSQLClientConfig() db.Config {
	host := "localhost"
	if os.Getenv("LOCAL_ENV") == "docker" {
		host = "mariadb"
	}

	return db.Config{
		User:         "root",
		Pass:         "",
		DBName:       "maria",
		Net:          "tcp",
		Host:         host,
		Port:         "3306",
		MaxOpenConns: 5,
		MaxIdleConns: 2,
	}
}
