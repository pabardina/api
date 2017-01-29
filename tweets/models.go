package tweets

import "github.com/jinzhu/gorm"

type Keyword struct {
	gorm.Model
	Label  string  `json:"label" gorm:"not null;unique"`
	Tweets []Tweet `json:"-"`
}

type Tweet struct {
	gorm.Model
	TweetID   string `json:"tweet_id" gorm:"not null;unique"`
	Likes     int    `json:"likes" gorm:"not null"`
	Retweets  int    `json:"retweets" gorm:"not null"`
	KeywordID uint   `json:"keyword_id" gorm:"not null"`
	Keyword   Keyword
}
