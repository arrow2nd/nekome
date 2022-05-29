package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// SearchRecentTweets ツイートを検索
func (a *API) SearchRecentTweets(query string, results int) ([]*twitter.TweetObj, error) {
	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.TweetRecentSearchOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	result, err := client.TweetRecentSearch(context.Background(), query, opts)
	if err != nil {
		return nil, fmt.Errorf("tweet search error: %v", err)
	}

	return result.Raw.Tweets, nil
}
