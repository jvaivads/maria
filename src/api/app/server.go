package main

import (
	"maria/src/api/db"
	"maria/src/api/user"
	"os"

	"github.com/gin-gonic/gin"
)

type controller interface {
	SetURLMapping(router *gin.Engine)
}

func main() {
	router := gin.Default()
	controllers := make([]controller, 0)

	controllers = append(controllers, user.NewController(
		user.NewService(
			user.NewRelationalDB(
				db.NewSQLClient(
					getSQLClientConfig(),
				)))))

	for i := range controllers {
		controllers[i].SetURLMapping(router)
	}

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
