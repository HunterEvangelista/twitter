package db

import (
	"testing"
)

func TestTweets(t *testing.T) {
	// test that we can get tweets
	tweets, err := GetTweets()
	if err != nil {
		t.Errorf("Error getting tweets: %v", err)
	}
	if len(tweets) == 0 {
		t.Errorf("Expected at least one tweet, got 0")
	}
}
