package mysql

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Config struct {
	Address            string
	Database           string
	User               string
	Password           string
	Timeout            time.Duration
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration
}

func NewSQLDatabase(config Config) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.Addr = config.Address
	cfg.User = config.User
	cfg.Passwd = config.Password
	cfg.DBName = config.Database
	cfg.Timeout = config.Timeout
	cfg.ReadTimeout = cfg.Timeout
	cfg.WriteTimeout = cfg.Timeout
	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(config.MaxOpenConnections)
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	return db, nil
}
