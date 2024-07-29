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

/*type Tweet struct {
	// will need to add id for author
	Id           int
	Author       string
	Content      string
	Likes        []string
	Favorites    []string
	Interestings []string
	PostDate     time.Time
}

func (t *Tweet) GetDate() string {
	return t.PostDate.Format("January 2, 2006")
}

// temporary id for tweets
var id = 0

func NewTweet(author, content string, likes, favorites, interestings []string, PostDate time.Time) *Tweet {
	id++
	return &Tweet{
		Id:           id,
		Author:       author,
		Content:      content,
		Likes:        likes,
		Favorites:    favorites,
		Interestings: interestings,
		PostDate:     PostDate,
	}
}
*/
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

	e.GET("/", func(c echo.Context) error {
		twts, err := DB.GetTweets()
		log.Println("Tweets: ", twts)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.Render(http.StatusOK, "index", twts)
	})
	// new route to get new tweet form
	//e.GET("/new-post", func(c echo.Context) error {
	//	return c.Render(http.StatusOK, "new-post", nil)
	//})
	//// new route to post new tweet and see it at the top of the feed
	//e.POST("/new-post", func(c echo.Context) error {
	//	author := "User"
	//	content := c.FormValue("tweet")
	//	likes := []string{}
	//	favorites := []string{}
	//	interestings := []string{}
	//	postDate := time.Now()
	//	newTweet := NewTweet(author, content, likes, favorites, interestings, postDate)
	//	data = append(Tweets{newTweet}, data...)
	//	return c.Render(http.StatusOK, "home", data)
	//})
	//
	//e.DELETE("/delete/:id", func(c echo.Context) error {
	//	id, err := strconv.Atoi(c.Param("id"))
	//	err = data.DeleteTweet(id)
	//	if err != nil {
	//		return c.JSON(http.StatusNotFound, err.Error())
	//	}
	//	log.Println("End work")
	//	return c.NoContent(200)
	//})
	//
	if err := e.Start(":9000"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
