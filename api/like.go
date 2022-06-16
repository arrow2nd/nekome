package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchLikedTweets ユーザのいいねしたツイートを取得
func (a *API) FetchLikedTweets(userID string, maxResults int) ([]*twitter.TweetObj, error) {
	opts := twitter.UserLikesLookupOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: maxResults,
	}

	result, err := a.client.UserLikesLookup(context.Background(), userID, opts)
	if err != nil {
		return nil, err
	}

	return result.Raw.Tweets, nil
}

// Like いいね
func (a *API) Like(tweetID string) error {
	if _, err := a.client.UserLikes(context.Background(), a.CurrentUser.ID, tweetID); err != nil {
		return fmt.Errorf("like tweet error: %v", err)
	}

	return nil
}

// UnLike いいねを解除
func (a *API) UnLike(tweetID string) error {
	if _, err := a.client.DeleteUserLikes(context.Background(), a.CurrentUser.ID, tweetID); err != nil {
		return fmt.Errorf("unlike tweet error: %v", err)
	}

	return nil
}
