package main

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

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

type Tweet struct {
	// will need to add id for author
	Author       string
	Content      string
	Likes        []string
	Favorites    []string
	Interestings []string
	PostDate     time.Time
}

func NewTweet(author, content string, likes, favorites, interestings []string, PostDate time.Time) *Tweet {
	return &Tweet{
		Author:       author,
		Content:      content,
		Likes:        likes,
		Favorites:    favorites,
		Interestings: interestings,
		PostDate:     PostDate,
	}
}

type Tweets []*Tweet

func main() {
	// will need to add env info
	// will need to connect to mongodb

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplate()
	e.Static("/", "public")

	// test data
	data := Tweets{
		NewTweet(
			"Jack Dorsey",
			"just setting up my twttr",
			[]string{"Elon Musk"},
			[]string{"Joe Biden"},
			[]string{"Kanye West"},
			time.Now(),
		),
		NewTweet(
			"Britney Spears",
			"Does anyone think global warming is a good thing?"+
				"I love Lady Gaga. I think she's a really interesting artist.",
			[]string{},
			[]string{},
			[]string{"Joe Biden", "Donald Trump", "Greta Thurnberg"},
			time.Date(2011, time.February, 1, 0, 0, 0, 0, time.UTC),
		),
		NewTweet(
			"Kevin Durant",
			"I'm watching the History channel in the club and I'm wondering"+
				"how do these people kno what's goin on on the sun"+
				"...ain't nobody ever been",
			[]string{"Lebron James", "Steph Curry"},
			[]string{"Kobe Bryant"},
			[]string{"Barack Obama", "Lary Bird"},
			time.Date(2010, time.July, 30, 0, 0, 0, 0, time.UTC),
		),
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})

	if err := e.Start(":9000"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
