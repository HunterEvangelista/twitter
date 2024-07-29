package db

import (
	"testing"
	"time"
)

func TestTweets(t *testing.T) {
	// test that we can get tweets
	DB, _ := NewDB()
	tweets, err := DB.GetTweets()
	if err != nil {
		t.Errorf("Error getting tweets: %v", err)
	}
	if len(tweets) == 0 {
		t.Errorf("Expected at least one tweet, got 0")
	}
}

func TestDB_CreateTweet(t *testing.T) {
	DB, _ := NewDB()
	tweet := &Tweet{
		Id:           0,
		Author:       "Default User",
		Content:      "This is a test tweet",
		Likes:        0,
		Favorites:    0,
		Interestings: 0,
		PostDate:     time.Now(),
	}
	err := DB.CreateTweet(tweet)
	if err != nil {
		t.Errorf("Error creating tweet: %v", err)
	}
}
