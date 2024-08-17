package main

import (
	"bytes"
	"encoding/json"
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
	CreatedDate time.Time
	Name        string
	Email       string
	Password    string
	UserId      int `json:"userId" form:"userId" query:"userId"`
}

type SignUpRequest struct {
	Name              string `json:"name" form:"name" query:"name"`
	Email             string `json:"email" form:"email" query:"email"`
	Password          string `json:"password" form:"password" query:"password"`
	ConfirmedPassword string `json:"confirmedPassword" form:"confirmedPassword" query:"confirmedPassword"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}

type ReadSessionRequest struct {
	Cookie http.Cookie `json:"cookie"`
}

type Data struct {
	User        *User
	ResponseMsg *ResponseMsg
	Tweets      db.Tweets
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
		Tweets:      db.Tweets{},
		User:        &User{},
		ResponseMsg: &ResponseMsg{},
	}

	e.GET("/", func(c echo.Context) error {
		// clear flash messages
		SessionData.ResponseMsg = &ResponseMsg{}
		var err error

		r, _ := http.NewRequest(http.MethodGet, "http://localhost:8733/read-session", nil)
		cookie, _ := c.Cookie("session")
		if cookie != nil {
			log.Println("Cookie: ", cookie.Value)
			log.Println("Cookie expirations: ", cookie.Expires)
			r.AddCookie(cookie)
		} else {
			log.Println("No cookie")
		}

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		responseData, _ := io.ReadAll(response.Body)
		err = json.Unmarshal(responseData, &SessionData.User)
		if err != nil {
			return err
		}

		if SessionData.User.UserId == 0 {
			return c.Render(http.StatusOK, "index", SessionData)
		}

		log.Println("Session Before call: ", SessionData.User)
		BuildHomePage(&SessionData, DB)
		log.Println("Session After call: ", SessionData.User)
		return c.Render(http.StatusOK, "index", SessionData)
	})

	e.GET("/signup", func(c echo.Context) error {
		return c.Render(http.StatusOK, "signup", nil)
	})

	e.POST("/signup", func(c echo.Context) error {
		var signupRequest SignUpRequest
		c.Bind(&signupRequest)
		log.Println(" signupRequest: ", signupRequest)
		signupJson, err := json.Marshal(signupRequest)
		body := bytes.NewReader(signupJson)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		r, _ := http.NewRequest(http.MethodPost, "http://localhost:6883/signup", body)
		r.Header.Add("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// parse response json into struct
		responseData, _ := io.ReadAll(response.Body)
		err = json.Unmarshal(responseData, &SessionData.ResponseMsg)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusOK {
			AssignFollowers(signupRequest.Email, DB)
			return c.Render(response.StatusCode, "login", SessionData)
		} else {
			return c.Render(http.StatusOK, "signup", SessionData)
		}
	})

	e.POST("/login", func(c echo.Context) error {
		// clear flash messages
		SessionData.ResponseMsg = &ResponseMsg{}
		log.Println("line 163")

		// get login request, bind to struct
		var loginRequest LoginRequest
		c.Bind(&loginRequest)
		log.Println("login request: ", loginRequest)
		loginJson, err := json.Marshal(loginRequest)
		if err != nil {
			return err
		}
		body := bytes.NewReader(loginJson)

		r, _ := http.NewRequest(http.MethodGet, "http://localhost:6883/authenticate", body)
		r.Header.Add("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// check with auth service if valid
		if response.StatusCode == http.StatusOK {
			// create session
			// get user id
			query := `
      SELECT Id FROM Users WHERE Email = ?
    `
			DB.QueryRow(query, loginRequest.Email).Scan(&SessionData.User.UserId)
			// create session with user id
			log.Println("User ID: ", SessionData.User.UserId)
			reqJson, _ := json.Marshal(SessionData.User)
			body := bytes.NewReader(reqJson)
			r, _ := http.NewRequest(http.MethodPost, "http://localhost:8733/create-session", body)
			r.Header.Add("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(r)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err.Error())
			}
			if resp.StatusCode != http.StatusOK {
				return c.JSON(http.StatusInternalServerError, "Error creating session")
			}
			// get cookie from response and store it
			c.SetCookie(resp.Cookies()[0])

			// flash error message
		} else {
			responseData, _ := io.ReadAll(response.Body)
			err = json.Unmarshal(responseData, &SessionData.ResponseMsg)
			if err != nil {
				return err
			}
			return c.Render(http.StatusOK, "login", SessionData)
		}
		// get user information
		query := `
    SELECT Name, Email, CreatedDate FROM Users WHERE Id = ?
    `

		err = DB.QueryRow(query, SessionData.User.UserId).Scan(&SessionData.User.Name, &SessionData.User.Email, &SessionData.User.CreatedDate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		BuildHomePage(&SessionData, DB)

		return c.Render(http.StatusOK, "index", SessionData)
	})

	e.POST("/logout", func(c echo.Context) error {
		// clear messages
		SessionData.ResponseMsg = &ResponseMsg{}
		// get cookie
		r, _ := http.NewRequest(http.MethodPost, "http://localhost:8733/delete-session", nil)
		cookie, _ := c.Cookie("session")
		if cookie != nil {
			r.AddCookie(cookie)
			cookie.Expires = time.Now().Add(-5 * time.Second)
			c.SetCookie(cookie)
		}

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		if resp.StatusCode == http.StatusOK {
			SessionData.User = &User{}
			SessionData.User.UserId = 0
			SessionData.ResponseMsg.Message = "Successfully logged out"
			SessionData.Tweets = db.Tweets{}
			return c.Render(http.StatusOK, "index", SessionData)
		} else {
			return c.JSON(http.StatusInternalServerError, "Error logging out")
		}
	})

	e.GET("/new-post", func(c echo.Context) error {
		return c.Render(http.StatusOK, "new-post", nil)
	})

	// new route to post new tweet and see it at the top of the feed
	e.POST("/new-post", func(c echo.Context) error {
		// clear flash messages
		SessionData.ResponseMsg = &ResponseMsg{}
		author := SessionData.User.Name
		content := c.FormValue("tweet")
		// check for profanity
		res := HandleProfanityCheck(content)
		if res != "" {
			SessionData.ResponseMsg.Message = res
			return c.Render(http.StatusOK, "new-post", SessionData)
		}
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

		err = DB.DeleteTweet(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(200)
	})

	e.POST("/like/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
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
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

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

func AssignFollowers(email string, db *db.DB) error {
	var newUserId int
	err := db.QueryRow("SELECT Id FROM Users WHERE Email = ?", email).Scan(&newUserId)
	if err != nil {
		return err
	}
	var userIds []int
	rows, err := db.Query("SELECT DISTINCT Id FROM Users WHERE Id != ?", newUserId)
	if err != nil {
		return err
	}

	for rows.Next() {
		var userId int
		err = rows.Scan(&userId)
		if err != nil {
			return err
		}
		userIds = append(userIds, userId)
	}

	// insert follow relationship
	for _, userId := range userIds {
		_, err = db.Exec("INSERT INTO Follows (FollowerId, FolloweeId) VALUES (?, ?)", newUserId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSessionUserInfo(userId int, db *db.DB) (*User, error) {
	var user User
	err := db.QueryRow("SELECT Name, Email, CreatedDate FROM Users WHERE Id = ?", userId).Scan(&user.Name, &user.Email, &user.CreatedDate)
	if err != nil {
		return nil, err
	}
	// add call to user info service
	userIdStr := strconv.Itoa(userId)
	http.NewRequest(http.MethodGet, "http://localhost:6050/get-user-info/"+userIdStr, nil)

	user.UserId = userId
	return &user, nil
}

type ProfanityResp struct {
	Profane string `json:"profane" form:"profane" query:"profane"`
}

func HandleProfanityCheck(content string) string {
	req := make(map[string]string)
	req["tweet"] = content
	log.Println("Content: ", content)
	body, _ := json.Marshal(req)
	reqBody := bytes.NewReader(body)
	r, _ := http.NewRequest(http.MethodGet, "http://localhost:6990/check", reqBody)
	r.Header.Add("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(r)
	responseData, _ := io.ReadAll(resp.Body)

	var result ProfanityResp
	err := json.Unmarshal(responseData, &result)
	log.Println(err)
	if err != nil {
		return "Error parsing response"
	}
	if result.Profane == "1" {
		return "Profanity detected"
	}
	return ""
}

func BuildHomePage(sessionData *Data, db *db.DB) error {
	var err error
	sessionData.User, _ = GetSessionUserInfo(sessionData.User.UserId, db)
	sessionData.Tweets, err = db.GetTweetsByFollowing(sessionData.User.UserId)
	if err != nil {
		return err
	}
	return nil
}
