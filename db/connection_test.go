package db

import (
	"database/sql"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	cfg := mysql.Config{
		User:   "root",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "twitter",
	}

	var err error
	conn, err = sql.Open("mysql", cfg.FormatDSN())
	assert.NoError(t, err, "Expected no error on sql.Open")

	pingErr := conn.Ping()
	assert.NoError(t, pingErr, "Expected no error on conn.Ping")

	if conn != nil {
		conn.Close()
	}
}
