package models

import (
	"github.com/jinzhu/gorm"
)

type Tweet struct {
	gorm.Model
	TweetID   string `json:"tweet_id,omitempty",gorm:"not null;unique"`
	Likes     int    `json:"likes,omitempty",gorm:"unique"`
	Retweets  int    `json:"retweets,omitempty",gorm:"unique"`
	KeywordID uint
	Keyword   Keyword
}

func GetTweetsByUser(db *gorm.DB) ([]Tweet, error) {
	tweets := []Tweet{}
	db.Preload("Keyword").Find(&tweets)
	return tweets, nil
}
