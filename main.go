package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"log"
	"main/db"
	"net/http"
	"strconv"
	"time"
)

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

//type Tweet struct {
//	// will need to add id for author
//	Id           int
//	Author       string
//	Content      string
//	Likes        []string
//	Favorites    []string
//	Interestings []string
//	PostDate     time.Time
//}
//
//func (t *Tweet) GetDate() string {
//	return t.PostDate.Format("January 2, 2006")
//}

func NewTweet(author, content string) *db.Tweet {
	return &db.Tweet{
		Id:           0,
		Author:       author,
		Content:      content,
		Likes:        0,
		Favorites:    0,
		Interestings: 0,
		PostDate:     time.Now(),
	}
}

//type Tweets []*Tweet
//
//func (t *Tweets) DeleteTweet(id int) error {
//	for i, tweet := range *t {
//		log.Println("Tweet id: ", tweet.Id)
//		log.Println("Selected id: ", id)
//		if tweet.Id == id {
//			log.Println("Tweet found")
//			*t = append((*t)[:i], (*t)[i+1:]...)
//			return nil
//		}
//	}
//	return errors.New("tweet not found")
//}

func main() {
	// will need to add env info
	// will need to connect to mongodb

	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplate()
	e.Static("/", "public")

	//// test data
	//data := Tweets{
	//	NewTweet(
	//		"Jack Dorsey",
	//		"just setting up my twttr",
	//		[]string{"Elon Musk"},
	//		[]string{"Joe Biden"},
	//		[]string{"Kanye West"},
	//		time.Now(),
	//	),
	//	NewTweet(
	//		"Britney Spears",
	//		"Does anyone think global warming is a good thing?"+
	//			"I love Lady Gaga. I think she's a really interesting artist.",
	//		[]string{},
	//		[]string{},
	//		[]string{"Joe Biden", "Donald Trump", "Greta Thunberg"},
	//		time.Date(2011, time.February, 1, 0, 0, 0, 0, time.UTC),
	//	),
	//	NewTweet(
	//		"Kevin Durant",
	//		"I'm watching the History channel in the club and I'm wondering"+
	//			"how do these people kno what's goin on on the sun"+
	//			"...ain't nobody ever been",
	//		[]string{"Lebron James", "Steph Curry"},
	//		[]string{"Kobe Bryant"},
	//		[]string{"Barack Obama", "Lary Bird"},
	//		time.Date(2010, time.July, 30, 0, 0, 0, 0, time.UTC),
	//	),
	//}
	DB, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	var twts db.Tweets

	e.GET("/", func(c echo.Context) error {
		var err error
		twts, err = DB.GetTweets()
		log.Println("Tweets: ", twts)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.Render(http.StatusOK, "index", twts)
	})

	e.GET("/new-post", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new-post", nil)
	})
	// new route to post new tweet and see it at the top of the feed
	e.POST("/new-post", func(c echo.Context) error {
		author := "Default User"
		content := c.FormValue("tweet")
		twt := NewTweet(author, content)
		err := DB.CreateTweet(twt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		err = DB.QueryRow("SELECT Id FROM Tweets WHERE Content = ?", content).Scan(&twt.Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		twts = append(db.Tweets{twt}, twts...)
		return c.Render(http.StatusOK, "home", twts)
	})

	e.DELETE("/delete/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = twts.DeleteTweet(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		err = DB.DeleteTweet(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(200)
	})

	e.POST("/like/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		// get current user id
		// just defaul user for now
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.LikeTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(200)
	})

	e.POST("/interesting/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		// get current user id
		// just default user for now
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.InterestingTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(200)
	})

	e.POST("/favorite/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		// get current user id
		// just default user for now
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.FavoriteTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(200)
	})

	if err := e.Start(":9000"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
