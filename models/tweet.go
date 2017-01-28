package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Tweet struct {
	gorm.Model
	TweetID   string `json:"tweet_id" gorm:"not null;unique"`
	Likes     int    `json:"likes" gorm:"not null"`
	Retweets  int    `json:"retweets" gorm:"not null"`
	KeywordID uint   `json:"keyword_id" gorm:"not null"`
	Keyword   Keyword
}

func GetTweetsByUser(db *gorm.DB) ([]Tweet, error) {
	tweets := []Tweet{}
	db.Preload("Keyword").Find(&tweets)
	return tweets, nil
}

func CreateTweet(db *gorm.DB, tweet *Tweet) error {

	if err := db.Create(&tweet).Error; err != nil {
		return err
	}
	db.Model(&tweet).Related(&Keyword{}, "KeywordID")
	db.Preload("Keyword").Model(&tweet)
	return nil
}

func ValidateTweet(tweet *Tweet) error {

	if tweet.TweetID == "" {
		return errors.New("Tweet_id must not be empty")
	}
	if tweet.Likes == 0 {
		return errors.New("Likes must not be empty")
	}
	if tweet.Retweets == 0 {
		return errors.New("Retweets must not be empty")
	}
	if tweet.KeywordID == 0 {
		return errors.New("Keyword must not be empty")
	}
	return nil
}
