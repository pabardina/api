package test

import (
	"github.com/hirondelle-app/api/tweets"
	"github.com/stretchr/testify/mock"
)

type MockTweetsManager struct {
	mock.Mock
}

func (m *MockTweetsManager) GetTweetsByUser() ([]tweets.Tweet, error) {
	args := m.Called()
	return args.Get(0).([]tweets.Tweet), args.Error(1)
}

func (m *MockTweetsManager) CreateTweet(tweet *tweets.Tweet) error {
	args := m.Called(tweet)
	return args.Error(0)
}

func (m *MockTweetsManager) ValidateTweet(tweet *tweets.Tweet) error {
	args := m.Called(tweet)
	return args.Error(0)
}

func (m *MockTweetsManager) CreateKeyword(keyword *tweets.Keyword) error {
	args := m.Called(keyword)
	return args.Error(0)
}

func (m *MockTweetsManager) GetKeywords() ([]tweets.Keyword, error) {
	args := m.Called()
	return args.Get(0).([]tweets.Keyword), args.Error(1)
}

func (m *MockTweetsManager) GetTweetsForKeyword(keywordID int, params *tweets.ParamsTweet) ([]tweets.Tweet, error) {
	args := m.Called(keywordID, params)
	return args.Get(0).([]tweets.Tweet), args.Error(1)
}

func (m *MockTweetsManager) GetCallsForMethod(methodName string) []mock.Call {
	calls := []mock.Call{}
	for _, call := range m.Calls {
		if call.Method == methodName {
			calls = append(calls, call)
		}
	}
	return calls
}
