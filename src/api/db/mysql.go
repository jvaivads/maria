package db

import (
	"database/sql"
	"sync"
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

	var err error
	if client, err = sql.Open("", ""); err != nil {
		panic(err)
	}

	client.SetMaxOpenConns(0)
	client.SetMaxIdleConns(0)
	client.SetConnMaxLifetime(0)
	client.SetConnMaxIdleTime(0)

	return client
}