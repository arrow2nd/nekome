package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// SearchRecentTweets : ツイートを検索
func (a *API) SearchRecentTweets(query, sinceId string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.TweetRecentSearchOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceId,
	}

	res, err := a.client.TweetRecentSearch(context.Background(), query, opts)
	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if res.Raw == nil {
		return []*twitter.TweetDictionary{}, res.RateLimit, nil
	}

	ok, tweets := createTweetSlice(res.Raw)
	if e := checkPartialError(res.Raw.Errors); !ok && e != nil {
		return nil, nil, e
	}

	return tweets, res.RateLimit, nil
}
