package api_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/hirondelle-app/api/api"
	. "github.com/hirondelle-app/api/common/test"
	"github.com/hirondelle-app/api/tweets"
	. "github.com/hirondelle-app/api/tweets/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Tweets", func() {

	var (
		request           *http.Request
		responseRecorder  *httptest.ResponseRecorder
		mockTweetsManager *MockTweetsManager
		tweetsHandlers    *TweetsHandlers
		err               error
	)

	BeforeEach(func() {
		responseRecorder = httptest.NewRecorder()
		mockTweetsManager = &MockTweetsManager{}
		tweetsHandlers = &TweetsHandlers{
			Manager: mockTweetsManager,
		}
	})

	Describe("GetTweetsEndpoint", func() {

		JustBeforeEach(func() {
			request, _ = http.NewRequest("GET", "/tweets", nil)
			tweetsHandlers.GetTweetsEndpoint(responseRecorder, request)
		})

		AfterEach(func() {
			err = nil
		})

		Context("when the manager successfully returns the tweets", func() {

			BeforeEach(func() {
				mockTweetsManager.On("GetAllTweets").Return([]tweets.Tweet{
					tweets.Tweet{TweetID: "ec815d46-e647-11e6-9902-8bf35f54ad22", Likes: 34, Retweets: 45, KeywordID: 56, Keyword: tweets.Keyword{Label: "python"}},
					tweets.Tweet{TweetID: "f0143618-e647-11e6-9656-4f834a284cc4", Likes: 43, Retweets: 76, KeywordID: 787, Keyword: tweets.Keyword{Label: "golang"}},
					tweets.Tweet{TweetID: "f53feefc-e647-11e6-8213-07d71760f8c3", Likes: 56, Retweets: 93, KeywordID: 12, Keyword: tweets.Keyword{Label: "java"}},
				}, nil)
			})

			It("should respond with a http status 200", func() {
				Expect(responseRecorder.Code).To(Equal(200))
			})

			It(`should respond with a content type "application/json"`, func() {
				Expect(responseRecorder.Header().Get("Content-Type")).To(Equal("application/json"))
			})

			It("should respond with a valid JSON", func() {
				Expect(responseRecorder.Body.String()).To(MatchJSON(ReadContentFileString("test/tweets-response-success.json")))
			})

			It("should get the tweet list from the underlying service", func() {
				Expect(len(mockTweetsManager.Calls)).To(Equal(1))
			})

		})

		Context("when the manager fails to return the tweets", func() {

			BeforeEach(func() {
				err = errors.New("Something really terrible happened!")
				mockTweetsManager.On("GetAllTweets").Return([]tweets.Tweet{}, err)
			})

			It("should respond with an invalid JSON", func() {
				Expect(responseRecorder.Body.String()).NotTo(MatchJSON(ReadContentFileString("test/tweets-response-success.json")))
			})

			It("Should respond with an error", func() {
				Expect(responseRecorder.Body.String()).To(MatchRegexp(err.Error()))
			})

		})

	})

	Describe("DeleteTweetEndpoint", func() {

		JustBeforeEach(func() {
			request, _ = http.NewRequest("DELETE", "/tweets/1", nil)
			m := mux.NewRouter()
			m.HandleFunc("/tweets/{tweetID}", tweetsHandlers.DeleteTweetEndpoint)
			m.ServeHTTP(responseRecorder, request)
		})

		AfterEach(func() {
			err = nil
		})

		Context("when the manager successfully delete the tweet", func() {

			BeforeEach(func() {
				mockTweetsManager.On("GetTweetByID", mock.Anything).Return(tweets.Tweet{
					TweetID:   "f53feefc-e647-11e6-8213-07d71760f8c3",
					Likes:     56,
					Retweets:  93,
					KeywordID: 12,
					Keyword:   tweets.Keyword{Label: "java"}},
					nil)
				mockTweetsManager.On("DeleteTweet", mock.Anything).Return(nil)
			})

			It("should respond with a http status 204", func() {
				Expect(responseRecorder.Code).To(Equal(204))
			})

		})

		Context("when the manager fails to find the tweet", func() {

			BeforeEach(func() {
				err = errors.New("There is an error !")
				mockTweetsManager.On("GetTweetByID", mock.Anything).Return(tweets.Tweet{}, err)
				mockTweetsManager.On("DeleteTweet", mock.Anything).Return(err)
			})

			It("should respond with a http status code 404", func() {
				Expect(responseRecorder.Code).To(Equal(404))
			})

			It("Should respond with an error", func() {
				Expect(responseRecorder.Body.String()).To(MatchRegexp(err.Error()))
			})

			It(`should respond with a content type "application/json"`, func() {
				Expect(responseRecorder.Header().Get("Content-Type")).To(Equal("application/json"))
			})

		})

		Context("when the manager fails to delete the tweet", func() {

			BeforeEach(func() {
				err = errors.New("There is an error !")
				mockTweetsManager.On("GetTweetByID", mock.Anything).Return(tweets.Tweet{
					TweetID:   "f53feefc-e647-11e6-8213-07d71760f8c3",
					Likes:     56,
					Retweets:  93,
					KeywordID: 12,
					Keyword:   tweets.Keyword{Label: "java"}},
					nil)
				mockTweetsManager.On("DeleteTweet", mock.Anything).Return(err)
			})

			It("should respond with a http status code 400", func() {
				Expect(responseRecorder.Code).To(Equal(400))
			})

			It("Should respond with an error", func() {
				Expect(responseRecorder.Body.String()).To(MatchRegexp(err.Error()))
			})

			It(`should respond with a content type "application/json"`, func() {
				Expect(responseRecorder.Header().Get("Content-Type")).To(Equal("application/json"))
			})

		})

	})
	// Describe("PostTweetEndpoint", func() {

	// 	JustBeforeEach(func() {
	// 		request, _ = http.NewRequest("POST", "/tweets", strings.NewReader(ReadContentFileString("test/tweet-to-create.json")))
	// 		tweetsHandlers.PostTweetEndpoint(responseRecorder, request)
	// 	})

	// 	Context("when the tweet is successfully created", func() {

	// 		BeforeEach(func() {
	// 			mockTweetsManager.On("ValidateTweet", mock.Anything).Return(nil)
	// 			mockTweetsManager.On("CreateTweet", mock.Anything).Return(nil)
	// 		})

	// 		It("should respond with a http status 201", func() {
	// 			Expect(responseRecorder.Code).To(Equal(201))
	// 		})

	// 		It("should validate the right tweet", func() {
	// 			validatedTweet := mockTweetsManager.GetCallsForMethod("ValidateTweet")[0].Arguments.Get(0).(*tweets.Tweet)
	// 			Expect(validatedTweet).To(Equal(&tweets.Tweet{
	// 				TweetID:   "ec815d46-e647-11e6-9902-8bf35f54ad22",
	// 				Likes:     34,
	// 				Retweets:  45,
	// 				KeywordID: 56,
	// 			}))
	// 		})

	// 		It("should create the right tweet", func() {
	// 			tweetUsedForCreation := mockTweetsManager.GetCallsForMethod("CreateTweet")[0].Arguments.Get(0).(*tweets.Tweet)
	// 			Expect(tweetUsedForCreation).To(Equal(&tweets.Tweet{
	// 				TweetID:   "ec815d46-e647-11e6-9902-8bf35f54ad22",
	// 				Likes:     34,
	// 				Retweets:  45,
	// 				KeywordID: 56,
	// 			}))
	// 		})

	// 	})

	// 	Context("when the tweet is invalid", func() {

	// 		var errorMsg string

	// 		BeforeEach(func() {
	// 			errorMsg = "The tweet is invalid !"
	// 			mockTweetsManager.On("ValidateTweet", mock.Anything).Return(errors.New(errorMsg))
	// 		})

	// 		It("should respond with a http status 400", func() {
	// 			Expect(responseRecorder.Code).To(Equal(400))
	// 		})

	// 		It("should respond with the incorrect message", func() {
	// 			listErrors := make(map[string]string)
	// 			listErrors["error"] = "invalid_tweet"
	// 			listErrors["error_description"] = errorMsg
	// 			jsonErrors, _ := json.Marshal(listErrors)
	// 			Expect(responseRecorder.Body.String()).To(MatchJSON(jsonErrors))
	// 		})
	// 	})

	// 	Context("when the tweet is not inserted in database", func() {

	// 		var errorMsg string

	// 		BeforeEach(func() {
	// 			mockTweetsManager.On("ValidateTweet", mock.Anything).Return(nil)

	// 			errorMsg = "There is an error with the database"
	// 			mockTweetsManager.On("CreateTweet", mock.Anything).Return(errors.New(errorMsg))

	// 		})

	// 		It("should respond with a http status 400", func() {
	// 			Expect(responseRecorder.Code).To(Equal(400))
	// 		})

	// 		It("should respond with the correct message", func() {
	// 			listErrors := make(map[string]string)
	// 			listErrors["error"] = "db_error"
	// 			listErrors["error_description"] = errorMsg
	// 			jsonErrors, _ := json.Marshal(listErrors)

	// 			Expect(responseRecorder.Body.String()).To(MatchJSON(jsonErrors))
	// 		})
	// 	})

	// })

})
