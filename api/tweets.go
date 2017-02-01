package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hirondelle-app/api/tweets"
)

type TweetsHandlers struct {
	Manager interface {
		GetTweetsByUser() ([]tweets.Tweet, error)
		CreateTweet(tweet *tweets.Tweet) error
		ValidateTweet(tweet *tweets.Tweet) error
		CreateKeyword(keyword *tweets.Keyword) error
		GetKeywords() ([]tweets.Keyword, error)
	} `inject:""`
}

func (h *TweetsHandlers) GetTweetsEndpoint(w http.ResponseWriter, req *http.Request) {
	//ctx := req.Context()
	//user := ctx.Value("user")

	tweetList, _ := h.Manager.GetTweetsByUser()

	if err := writeJSON(w, tweetList, 200); err != nil {
		log.Fatal(err)
	}
}

func (h *TweetsHandlers) PostTweetEndpoint(w http.ResponseWriter, req *http.Request) {
	tweet := tweets.Tweet{}
	json.NewDecoder(req.Body).Decode(&tweet)

	// check if tweet is correct
	if err := h.Manager.ValidateTweet(&tweet); err != nil {
		httpError(w, 400, "invalid_tweet", err.Error())
		return
	}

	// save tweet
	if err := h.Manager.CreateTweet(&tweet); err != nil {
		httpError(w, 400, "db_error", err.Error())
		return
	}

	// return tweet in response
	writeJSON(w, tweet, 201)
}

func (h *TweetsHandlers) PostKeywordEndpoint(w http.ResponseWriter, req *http.Request) {

	keyword := tweets.Keyword{}
	json.NewDecoder(req.Body).Decode(&keyword)

	if keyword.Label == "" {
		httpError(w, 400, "invalid_keyword", "Label must not be empty")
		return
	}

	if err := h.Manager.CreateKeyword(&keyword); err != nil {
		httpError(w, 400, "db_error", err.Error())
		return
	}

	writeJSON(w, keyword, 201)
}

func (h *TweetsHandlers) GetAllKeywordsEndpoint(w http.ResponseWriter, req *http.Request) {

	keywords, _ := h.Manager.GetKeywords()

	if err := writeJSON(w, keywords, 200); err != nil {
		log.Fatal(err)
	}
}
