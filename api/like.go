package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// UserLikesLookup ユーザのいいねしたツイートを取得
func (a *API) UserLikesLookup(userID string, maxResults int) ([]*twitter.TweetObj, error) {
	client := a.newClient(a.CurrentUser.Token)

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

// Like いいね
func (a *API) Like(tweetID string) error {
	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.UserLikes(context.Background(), a.CurrentUser.ID, tweetID); err != nil {
		return fmt.Errorf("like tweet error: %v", err)
	}

	return nil
}
