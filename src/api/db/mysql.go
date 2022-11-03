package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/go-sql-driver/mysql"
)

var (
	client *sql.DB
	mu     sync.Mutex
)

func GetSQLClient() Client {
	mu.Lock()
	defer mu.Unlock()

	if client != nil {
		return client
	}

	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Addr = "localhost:3306"
	cfg.DBName = "maria"

	var err error
	if client, err = sql.Open("mysql", cfg.FormatDSN()); err != nil {
		panic(err)
	}

	/*
		client.SetMaxOpenConns(0)
		client.SetMaxIdleConns(0)
		client.SetConnMaxLifetime(0)
		client.SetConnMaxIdleTime(0)
	*/

	if err = client.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("database connected")

	return client
}
