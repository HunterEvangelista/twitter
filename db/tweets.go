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
	var tweets Tweets
	for rows.Next() {
		var tweet Tweet
		err = rows.Scan(&tweet.Id, &tweet.Content, &tweet.PostDate, &tweet.Author)
		if err != nil {
			return nil, fmt.Errorf("error scanning tweets: %v", err)
		}
		// find likes, interestings, and favorites
		likesQuery := `
		SELECT COUNT(*) as likes FROM Likes WHERE TweetId = ?
		`
		err = db.QueryRow(likesQuery, tweet.Id).Scan(&tweet.Likes)

		interestingsQuery := `
        		SELECT COUNT(*) as interestings FROM Interestings WHERE TweetId = ?
		`
		err = db.QueryRow(interestingsQuery, tweet.Id).Scan(&tweet.Interestings)

		favoritesQuery := `
				SELECT COUNT(*) as favorites FROM Favorites WHERE TweetId = ?
		`
		err = db.QueryRow(favoritesQuery, tweet.Id).Scan(&tweet.Favorites)
		tweets = append(tweets, &tweet)
	}
	return tweets, nil
}

func (db *DB) CreateTweet(tweet *Tweet) error {
	query := `
	INSERT INTO Tweets (Content, PostDate, UserId)
	VALUES (?, ?, (SELECT Id FROM Users WHERE Name = ?))
	`
	_, err := db.Exec(query, tweet.Content, tweet.PostDate.Format("2006-01-01 15:04:06"), tweet.Author)
	if err != nil {
		return fmt.Errorf("error inserting tweet: %v", err)
	}
	return nil
}

func (db *DB) DeleteTweet(id int) error {
	query := `
	DELETE FROM Tweets WHERE Id = ?
	`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting tweet: %v", err)
	}
	return nil
}

// this should be moved to a controller eventually
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

func (db *DB) LikeTweet(tweetId, userId int) error {
	query := `
	INSERT INTO Likes (TweetId, UserId)
	VALUES (?, ?)
	`
	_, err := db.Exec(query, tweetId, userId)
	if err != nil {
		return fmt.Errorf("error liking tweet: %v", err)
	}
	return nil
}

func (db *DB) InterestingTweet(tweetId, userId int) error {
	query := `
	INSERT INTO Interestings (TweetId, UserId)
	VALUES (?, ?)
	`
	_, err := db.Exec(query, tweetId, userId)
	if err != nil {
		return fmt.Errorf("error interesting tweet: %v", err)
	}
	return nil
}

func (db *DB) FavoriteTweet(tweetId, userId int) error {
	query := `
	INSERT INTO Favorites (TweetId, UserId)
	VALUES (?, ?)
	`
	_, err := db.Exec(query, tweetId, userId)
	if err != nil {
		return fmt.Errorf("error favoriting tweet: %v", err)
	}
	return nil
}
