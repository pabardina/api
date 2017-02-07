package tweets

import "fmt"

type PaginateTweet struct {
	Start    int     `json:"start"`
	Limit    int     `json:"limit"`
	Total    int     `json:"total"`
	Next     string  `json:"next"`
	Previous string  `json:"previous"`
	Results  []Tweet `json:"results"`
}

type ParamsTweet struct {
	Limit     int
	Start     int
	Retweets  int
	Likes     int
	KeywordID int
	Total     int
}

func GetTweetsPagination(tweets []Tweet, params *ParamsTweet) (PaginateTweet, error) {

	var previous, next string

	if params.Start >= 1 {
		previousVal := 0
		previous = fmt.Sprintf("/keywords/%d/tweets?start=%d&limit=%d&retweets=%d&likes=%d",
			params.KeywordID, previousVal, params.Limit, params.Retweets, params.Likes)
	}

	if (params.Start + params.Limit) <= params.Total {
		nextVal := params.Start + params.Limit
		next = fmt.Sprintf("/keywords/%d/tweets?start=%d&limit=%d&retweets=%d&likes=%d",
			params.KeywordID, nextVal, params.Limit, params.Retweets, params.Likes)
	}

	paginateTweet := PaginateTweet{
		Start:    params.Start,
		Limit:    params.Limit,
		Total:    params.Total,
		Next:     next,
		Previous: previous,
		Results:  tweets,
	}

	return paginateTweet, nil
}
