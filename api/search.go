package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// SearchRecentTweets : ツイートを検索
func (a *API) SearchRecentTweets(query, sinceID string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.TweetRecentSearchOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.TweetRecentSearch(context.Background(), query, opts)

	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if e := checkPartialError(res.Raw.Errors); e != nil {
		return nil, nil, e
	}

	return createTweetDictionarySlice(res.Raw), res.RateLimit, nil
}
