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
	Name              string `json:"name" form:"name" query:"name"`
	Email             string `json:"email" form:"email" query:"email"`
	Password          string `json:"password" form:"password" query:"password"`
	ConfirmedPassword string `json:"confirmedPassword" form:"confirmedPassword" query:"confirmedPassword"`
}

type ResponseMessage struct {
	Message string `json:"message"`
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

		log.Println(userSignUpRequest)

		if userSignUpRequest.Password != userSignUpRequest.ConfirmedPassword {
			return c.JSON(400, ResponseMessage{Message: "Passwords do not match"})
		}
		log.Println("Passwords match")
		// check if email already exists
		query := `
    SELECT COUNT(*) FROM Users WHERE Email = ?
    `
		var count int
		err := db.QueryRow(query, userSignUpRequest.Email).Scan(&count)
		if err != nil {
			return c.JSON(500, ResponseMessage{Message: "Internal server error"})
		}
		log.Println("Count: ", count)
		if count != 0 {
			return c.JSON(400, ResponseMessage{Message: "Email already exists"})
		}
		log.Println("Email does not exist")

		// insert user
		query = `
    INSERT INTO Users (Name, Email, Password, CreatedDate) VALUES (?, ?, ?, NOW()) 
    `
		_, err = db.Exec(query, userSignUpRequest.Name, userSignUpRequest.Email, userSignUpRequest.Password)
		if err != nil {
			return c.JSON(500, ResponseMessage{Message: "Internal server error"})
		}

		log.Println("User created")

		return c.JSON(200, ResponseMessage{Message: "Account Created, please login."})
	})

	if err := e.Start(":6883"); err != nil {
		log.Fatal(err)
	}
}
