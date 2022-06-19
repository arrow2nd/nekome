package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// SearchRecentTweets ツイートを検索
func (a *API) SearchRecentTweets(query, sinceID string, results int) ([]*twitter.TweetDictionary, error) {
	opts := twitter.TweetRecentSearchOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.TweetRecentSearch(context.Background(), query, opts)
	if err != nil {
		return nil, fmt.Errorf("tweet search error: %v", err)
	}

	return createTweetDictionarySlice(res.Raw), nil
}
