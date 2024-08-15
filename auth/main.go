package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// looking to enable basic http auth
// get email and password as form values
// note that emails are guarenteed unique

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewDB() (*sql.DB, error) {
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
	return conn, nil
}

var db *sql.DB

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	db, _ = NewDB()

	e.GET("/authenticate", func(c echo.Context) error {
		var loginRequest LoginRequest
		if err := c.Bind(&loginRequest); err != nil {
			return err
		}

		query := `
      SELECT Password FROM Users WHERE Email = ?
    `

		var password string
		err := db.QueryRow(query, loginRequest.Email).Scan(&password)
		if err != nil {
			return c.String(404, "Email not found")
		}

		if password != loginRequest.Password {
			return c.String(401, "Incorrect password")
		} else {
			return c.String(200, "Authorized")
		}
	})

	if err := e.Start(":6883"); err != nil {
		log.Fatal(err)
	}
}
