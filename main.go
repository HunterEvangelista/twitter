package main

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// tmeplate info
type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	// will need to add env info
	// will need to connect to mongodb

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplate()
	e.Static("/", "public")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	if err := e.Start(":9000"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
