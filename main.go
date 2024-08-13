package main

import (
	"errors"
	"html/template"
	"io"
	"log"
	"main/db"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

type User struct {
	UserId      int
	Name        string
	Email       string
	Password    string
	CreatedDate time.Time
}

type Data struct {
	Tweets db.Tweets
	User   *User
}

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

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplate()
	e.Static("/", "public")

	DB, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	SessionData := Data{
		Tweets: db.Tweets{},
		User:   &User{},
	}

	DefaultUser := &User{
		UserId:      6,
		Name:        "Default User",
		Email:       "dev@test.com",
		Password:    "password",
		CreatedDate: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	SessionData.User = DefaultUser

	e.GET("/", func(c echo.Context) error {
		var err error
		SessionData.Tweets, err = DB.GetTweets()
		log.Println("should be true: ", SessionData.Tweets[0].IsLiked(6))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.Render(http.StatusOK, "index", SessionData)
	})

	e.GET("/new-post", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new-post", nil)
	})
	// new route to post new tweet and see it at the top of the feed
	e.POST("/new-post", func(c echo.Context) error {
		author := SessionData.User.Name
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
		SessionData.Tweets = append(db.Tweets{twt}, SessionData.Tweets...)
		return c.Render(http.StatusOK, "home", SessionData)
	})

	e.GET("expand/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		tweet, err := DB.GetTweet(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		log.Println("GetDate: ", tweet.Likes)
		return c.Render(http.StatusOK,
			"expanded-tweet",
			tweet)
	})

	e.DELETE("/delete/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		log.Println("Tweets: ", SessionData.Tweets)
		SessionData.Tweets.DeleteTweet(id)
		log.Println("Tweets: ", SessionData.Tweets)

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
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.LikeTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		likedTweet := SessionData.Tweets.GetTweetById(id)
		likedTweet.Likes++
		return c.Render(200, "liked", likedTweet)
	})

	e.DELETE("/unlike/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.UnlikeTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		unlikedTweet := SessionData.Tweets.GetTweetById(id)
		unlikedTweet.Likes--
		return c.Render(200, "unliked", unlikedTweet)
	})

	e.POST("/interesting/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.InterestingTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		interestingTweet := SessionData.Tweets.GetTweetById(id)
		interestingTweet.Interestings++
		return c.Render(200, "interesting", interestingTweet)
	})

	e.DELETE("/uninteresting/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.UninterestingTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		uninterestingTweet := SessionData.Tweets.GetTweetById(id)
		uninterestingTweet.Interestings--
		return c.Render(200, "uninteresting", uninterestingTweet)
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
		favoritedTweet := SessionData.Tweets.GetTweetById(id)
		favoritedTweet.Favorites++
		return c.Render(200, "favorited", favoritedTweet)
	})

	e.DELETE("/unfavorite/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		var userId int
		err = DB.QueryRow("SELECT Id FROM Users WHERE Name = ?", "Default User").Scan(&userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = DB.UnfavoriteTweet(id, userId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		unfavoritedTweet := SessionData.Tweets.GetTweetById(id)
		unfavoritedTweet.Favorites--
		return c.Render(200, "unfavorited", unfavoritedTweet)
	})

	if err := e.Start(":9000"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
