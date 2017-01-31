package tweets_test

import (
	. "github.com/hirondelle-app/api/tweets"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TweetsManager", func() {

	var manager *Manager

	BeforeEach(func() {
		manager = &Manager{}
	})

	Describe("ValidateTweet", func() {

		var validTweet, tweetToValidate Tweet
		var err error

		AfterEach(func() {
			err = nil
		})

		BeforeEach(func() {
			validTweet = Tweet{
				TweetID:   "549849ce-e640-11e6-8386-ff27b143351b",
				Likes:     13,
				Retweets:  45,
				KeywordID: 23,
			}
		})

		JustBeforeEach(func() {
			err = manager.ValidateTweet(&tweetToValidate)
		})

		Context("when everying is ok", func() {

			BeforeEach(func() {
				tweetToValidate = validTweet
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

		})

		Context("when the tweet does not have an id", func() {

			BeforeEach(func() {
				validTweet.TweetID = ""
				tweetToValidate = validTweet
			})

			It("should return the right error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Tweet_id must not be empty"))
			})

		})

		Context("when the tweet does not have likes", func() {

			BeforeEach(func() {
				validTweet.Likes = 0
				tweetToValidate = validTweet
			})

			It("should return the right error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Likes must not be empty"))
			})
		})

		Context("when the tweet does not have retweets", func() {

			BeforeEach(func() {
				validTweet.Retweets = 0
				tweetToValidate = validTweet
			})

			It("should return the right error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Retweets must not be empty"))
			})

		})

		Context("when the tweet does not have keyword ID", func() {

			BeforeEach(func() {
				validTweet.KeywordID = 0
				tweetToValidate = validTweet
			})

			It("should return the right error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("Keyword ID must not be empty"))
			})

		})
	})

})
