package tweets

type PaginateTweet struct {
	Start    int     `json:"start"`
	PerPage  int     `json:"perPage"`
	Total    int     `json:"total"`
	Next     string  `json:"next"`
	Previous string  `json:"previous"`
	Results  []Tweet `json:"results"`
}

type ParamsTweet struct {
	Limit    int
	Start    int
	Retweets int
	Likes    int
}
