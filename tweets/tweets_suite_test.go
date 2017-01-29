package tweets_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTweets(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tweets Suite")
}
