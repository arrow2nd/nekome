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

// AddBookmark ツイートをブックマークに追加
// FIXME: Invalid Request:One or more parameters to your request was invalid
func (a *API) AddBookmark(tweetID string) error {
	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.AddTweetBookmark(context.Background(), a.CurrentUser.ID, tweetID); err != nil {
		return fmt.Errorf("add bookmark: %v", err)
	}

	return nil
}

// DeleteBookmark ツイートをブックマークから削除
func (a *API) DeleteBookmark(tweetID string) error {
	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.RemoveTweetBookmark(context.Background(), a.CurrentUser.ID, tweetID); err != nil {
		return fmt.Errorf("delete bookmark: %v", err)
	}

	return nil
}
