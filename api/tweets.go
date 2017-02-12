package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hirondelle-app/api/tweets"
)

type TweetsHandlers struct {
	Manager interface {
		GetAllTweets() ([]tweets.Tweet, error)
		GetTweetByID(tweetID int) (tweets.Tweet, error)
		DeleteTweet(tweet *tweets.Tweet) error
		ValidateTweet(tweet *tweets.Tweet) error
		CreateKeyword(keyword *tweets.Keyword) error
		DeleteKeyword(keywordID int) error
		GetKeywordByID(keywordID int) (tweets.Keyword, error)
		GetKeywords() ([]tweets.Keyword, error)
		GetTweetsForKeyword(keywordID int, params *tweets.ParamsTweet) (tweets.PaginateTweet, error)
	} `inject:""`
}

func (h *TweetsHandlers) GetTweetsEndpoint(w http.ResponseWriter, req *http.Request) {
	tweetList, err := h.Manager.GetAllTweets()

	if err != nil {
		httpError(w, 400, "error", err.Error())
		return
	}

	if err := writeJSON(w, tweetList, 200); err != nil {
		log.Fatal(err)
	}
}

func (h *TweetsHandlers) DeleteTweetEndpoint(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	tweetStr := vars["tweetID"]

	if tweetStr == "" {
		httpError(w, 404, "invalid_tweet", "tweetID must be set in url")
		return
	}

	tweetID, _ := strconv.Atoi(tweetStr)

	tweet, err := h.Manager.GetTweetByID(tweetID)
	if err != nil {
		httpError(w, 404, "not_found", err.Error())
		return
	}

	if err := h.Manager.DeleteTweet(&tweet); err != nil {
		httpError(w, 400, "db_error", err.Error())
		return
	}

	writeJSON(w, nil, 204)
}

func (h *TweetsHandlers) PostKeywordEndpoint(w http.ResponseWriter, req *http.Request) {

	keyword := tweets.Keyword{}

	if err := json.NewDecoder(req.Body).Decode(&keyword); err != nil {
		httpError(w, 400, "invalid_keyword", "Label must not be empty")
		return
	}

	if err := h.Manager.CreateKeyword(&keyword); err != nil {
		httpError(w, 400, "db_error", err.Error())
		return
	}

	writeJSON(w, keyword, 201)
}

func (h *TweetsHandlers) DeleteKeywordEndpoint(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	keywordStr := vars["keywordID"]

	if keywordStr == "" {
		httpError(w, 404, "invalid_keyword", "keywordID must be set in url")
		return
	}

	keywordID, _ := strconv.Atoi(keywordStr)

	if err := h.Manager.DeleteKeyword(keywordID); err != nil {
		httpError(w, 400, "db_error", err.Error())
		return
	}

	writeJSON(w, nil, 204)
}

func (h *TweetsHandlers) GetAllKeywordsEndpoint(w http.ResponseWriter, req *http.Request) {

	keywords, _ := h.Manager.GetKeywords()

	if err := writeJSON(w, keywords, 200); err != nil {
		log.Fatal(err)
	}
}

func (h *TweetsHandlers) GetTweetsByKeywordEndpoint(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	keywordStr := vars["keywordID"]
	keywordID, _ := strconv.Atoi(keywordStr)

	// really need to check ?
	if keywordStr == "" {
		log.Fatal("empty")
	}

	params := &tweets.ParamsTweet{
		Limit:     GetQueryParamToStr("limit", req),
		Start:     GetQueryParamToStr("start", req),
		Retweets:  GetQueryParamToStr("retweets", req),
		Likes:     GetQueryParamToStr("likes", req),
		KeywordID: keywordID,
	}

	if params.Limit == 0 {
		params.Limit = 25
	}

	tweets, err := h.Manager.GetTweetsForKeyword(keywordID, params)
	if err != nil {
		log.Fatal(err)
	}

	if err := writeJSON(w, tweets, 200); err != nil {
		log.Fatal(err)
	}
}
