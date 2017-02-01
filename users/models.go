package users

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	TwitterID string `json:"twitter_id"`

	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
	RetweetLimit      int    `json:"retweet_limit"`
	LikeLimit         int    `json:"like_limit"`
	//ContainerID       string `json:"container_id"`
	//Keywords          []tweets.Keyword `json:"keywords" gorm:"many2many:user_keywords;"`
}
