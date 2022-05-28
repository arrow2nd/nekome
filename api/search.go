package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// TweetRecentSearch ツイートを検索
func (a *API) TweetRecentSearch(query string, results int) ([]*twitter.TweetObj, error) {
	client, err := a.newClient(a.CurrentUser.Token)
	if err != nil {
		return nil, err
	}

	opts := twitter.TweetRecentSearchOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	searchResults, err := client.TweetRecentSearch(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}

	return searchResults.Raw.Tweets, nil
}
