package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmedPassword"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	UserId  int    `json:"userId"`
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

	e.POST("signup", func(c echo.Context) error {
		var userSignUpRequest SignupRequest
		c.Bind(&userSignUpRequest)

		if userSignUpRequest.Password != userSignUpRequest.ConfirmedPassword {
			return c.JSON(400, ErrorResponse{Message: "Passwords do not match"})
		}
		// check if email already exists
		query := `
    SELECT COUNT(*) FROM Users WHERE Email = ?
    `
		var count int
		err := db.QueryRow(query, userSignUpRequest.Email).Scan(&count)
		if err != nil {
			return c.JSON(500, ErrorResponse{Message: "Internal server error"})
		}
		if count != 0 {
			return c.JSON(400, ErrorResponse{Message: "Email already exists"})
		}

		// insert user
		query = `
    INSERT INTO Users (Name, Email, Password, CreatedDate) VALUES (?, ?, ?, NOW()) 
    `
		_, err = db.Exec(query, userSignUpRequest.Name, userSignUpRequest.Email, userSignUpRequest.Password)
		if err != nil {
			return c.JSON(500, ErrorResponse{Message: "Internal server error"})
		}

		var successResponse SuccessResponse
		successResponse.Message = "User created"

		// get new user id
		query = `
      SELECT Id FROM Users WHERE Email = ?
    `
		err = db.QueryRow(query, userSignUpRequest.Email).Scan(&successResponse.UserId)
		if err != nil {
			return c.JSON(500, ErrorResponse{Message: "Internal server error"})
		}
		return c.JSON(200, successResponse)
	})

	if err := e.Start(":6883"); err != nil {
		log.Fatal(err)
	}
}
