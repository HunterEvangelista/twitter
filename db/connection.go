package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	cfg := mysql.Config{
		User:      "root",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "twitter",
		ParseTime: true,
	}
	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	return &DB{conn}, nil
}
