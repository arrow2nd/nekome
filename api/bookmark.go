package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchBookmarkTweets ブックマークしたツイートを取得
func (a *API) FetchBookmarkTweets() ([]*twitter.TweetObj, error) {
	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.TweetBookmarksLookupOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
	}

	result, err := client.TweetBookmarksLookup(context.Background(), a.CurrentUser.ID, opts)
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks error: %v", err)
	}

	return result.Raw.Tweets, nil
}
