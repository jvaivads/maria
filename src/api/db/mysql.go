package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	User   string
	Pass   string
	DBName string
	Net    string
	Host   string
	Port   string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
}

func (cfg Config) toMySQLConfig() *mysql.Config {
	mycfg := mysql.NewConfig()
	mycfg.User = cfg.User
	mycfg.Passwd = cfg.Pass
	mycfg.DBName = cfg.DBName
	mycfg.Net = cfg.Net
	mycfg.Addr = cfg.Host + ":" + cfg.Port

	return mycfg
}

func NewSQLClient(cfg Config) Client {
	var (
		client *sql.DB
		err    error
	)

	if client, err = sql.Open("mysql", cfg.toMySQLConfig().FormatDSN()); err != nil {
		panic(err)
	}

	client.SetMaxOpenConns(cfg.MaxOpenConns)
	client.SetMaxIdleConns(cfg.MaxIdleConns)

	if err = client.Ping(); err != nil {
		panic(fmt.Errorf("cannot connect with DSN %s due to: %w", cfg.toMySQLConfig().FormatDSN(), err))
	}

	fmt.Println("database connected")

	return client
}
