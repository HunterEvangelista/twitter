package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type SessionUserId struct {
	ID int `json:"userId"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("No .env file found")
	}
	e := echo.New()
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	e.Use(session.Middleware(store), middleware.Logger())

	e.GET(
		"/create-session",
		func(c echo.Context) error {
			sess, err := store.Get(c.Request(), "session")
			if err != nil {
				log.Println("Error getting session")
				return err
			}

			sess.Values["UserId"] = c.FormValue("userId")
			if err := sess.Save(c.Request(), c.Response()); err != nil {
				log.Println("Error Saving session")
				return err
			}
			return c.String(http.StatusOK, "Session created")
		})

	e.GET("/read-session", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		// access the session value
		sessId, ok := sess.Values["UserId"].(string)
		var sessIdInt int
		if !ok {
			sessIdInt = 0
		} else {
			sessIdInt, _ = strconv.Atoi(sessId)
		}
		userId := &SessionUserId{ID: sessIdInt}
		log.Println("Session ID: ", userId.ID)
		return c.JSON(http.StatusOK, userId)
	})

	e.GET("delete-session", func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		sess.Options.MaxAge = -1
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}
		return c.String(http.StatusOK, "Session deleted")
	})

	if err := e.Start(":8733"); err != nil {
		log.Fatal(err)
	}
}
