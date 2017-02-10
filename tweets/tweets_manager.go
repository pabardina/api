package tweets

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Manager struct {
	DB *gorm.DB `inject:""`
}

func (m *Manager) GetTweetsByUser() ([]Tweet, error) {
	tweetList := []Tweet{}
	m.DB.Preload("Keyword").Find(&tweetList)
	return tweetList, nil
}

func (m *Manager) CreateTweet(tweet *Tweet) error {
	if err := m.DB.Create(&tweet).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) ValidateTweet(tweet *Tweet) error {
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
		return errors.New("Keyword ID must not be empty")
	}
	return nil
}

func (m *Manager) CreateKeyword(keyword *Keyword) error {
	if err := m.DB.Create(&keyword).Error; err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetKeywords() ([]Keyword, error) {
	keywords := []Keyword{}
	m.DB.Preload("Tweets").Find(&keywords)
	return keywords, nil
}

func (m *Manager) GetTweetsForKeyword(keywordID int, params *ParamsTweet) (PaginateTweet, error) {

	tweets := []Tweet{}
	results := m.DB.Where("keyword_id = ? AND likes >= ? AND retweets >= ?", keywordID, params.Likes, params.Retweets).
		Order("created_at desc").Find(&tweets)
	params.Total = len(tweets)
	results.Offset(params.Start).Limit(params.Limit).Preload("Keyword").Find(&tweets)

	paginateTweet, err := GetTweetsPagination(tweets, params)

	if err != nil {
		return PaginateTweet{}, err
	}

	return paginateTweet, nil

}
