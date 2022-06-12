package api

import "github.com/g8rswimmer/go-twitter/v2"

func createTweetDictionarySlice(raw *twitter.TweetRaw) []*twitter.TweetDictionary {
	contents := []*twitter.TweetDictionary{}
	dics := raw.TweetDictionaries()

	for _, tweet := range raw.Tweets {
		contents = append(contents, dics[tweet.ID])
	}

	return contents
}