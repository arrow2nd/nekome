package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// UserLikesLookup ユーザのいいねしたツイートを取得
func (a *API) UserLikesLookup(userID string, maxResults int) ([]*twitter.TweetObj, error) {
	client, err := a.newClient(a.CurrentUser.Token)
	if err != nil {
		return nil, err
	}

	opts := twitter.UserLikesLookupOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: maxResults,
	}

	results, err := client.UserLikesLookup(context.Background(), userID, opts)
	if err != nil {
		return nil, err
	}

	return results.Raw.Tweets, nil
}
