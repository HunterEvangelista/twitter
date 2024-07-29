package db

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Tweet struct {
	Id           int
	Author       string
	Content      string
	Likes        int
	Favorites    int
	Interestings int
	PostDate     time.Time
}

func (t *Tweet) GetDate() string {
	return t.PostDate.Format("January 2, 2006")
}

//func NewTweet(id int, author, content string, likes, favorites, interestings int, PostDate time.Time) *Tweet {
//	return &Tweet{
//		Id:           id,
//		Author:       author,
//		Content:      content,
//		Likes:        likes,
//		Favorites:    favorites,
//		Interestings: interestings,
//		PostDate:     PostDate,
//	}
//}

type Tweets []*Tweet

func (db *DB) GetTweets() (Tweets, error) {
	query := `
	SELECT T.ID, T.Content, T.PostDate, U.Name
	FROM Tweets T
	INNER JOIN Users U ON T.UserId = U.Id
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting tweets: %v", err)
	}
	log.Println("made it past first query")
	var tweets Tweets
	for rows.Next() {
		var tweet Tweet
		err = rows.Scan(&tweet.Id, &tweet.Content, &tweet.PostDate, &tweet.Author)
		if err != nil {
			return nil, fmt.Errorf("error scanning tweets: %v", err)
		}
		log.Println("Got this tweet: ", tweet)
		// find likes, interestings, and favorites
		likesQuery := `
		SELECT COUNT(*) as likes FROM Likes WHERE TweetId = ?
		`
		err = db.QueryRow(likesQuery, tweet.Id).Scan(&tweet.Likes)
		log.Println("Likes: ", tweet.Likes)

		interestingsQuery := `
        		SELECT COUNT(*) as interestings FROM Interestings WHERE TweetId = ?
		`
		err = db.QueryRow(interestingsQuery, tweet.Id).Scan(&tweet.Interestings)
		log.Println("Interestings: ", tweet.Interestings)

		favoritesQuery := `
				SELECT COUNT(*) as favorites FROM Favorites WHERE TweetId = ?
		`
		err = db.QueryRow(favoritesQuery, tweet.Id).Scan(&tweet.Favorites)
		log.Println("Favorites: ", tweet.Favorites)
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}

func (t *Tweets) DeleteTweet(id int) error {
	for i, tweet := range *t {
		log.Println("Tweet id: ", tweet.Id)
		log.Println("Selected id: ", id)
		if tweet.Id == id {
			log.Println("Tweet found")
			*t = append((*t)[:i], (*t)[i+1:]...)
			return nil
		}
	}
	return errors.New("tweet not found")
}
