package api

import (
	"context"
)

// Retweet : リツイート
func (a *API) Retweet(tweetID string) error {
	_, err := a.client.UserRetweet(context.Background(), a.CurrentUser.ID, tweetID)

	return checkError(err)
}

// UnRetweet : リツイートを解除
func (a *API) UnRetweet(tweetID string) error {
	_, err := a.client.DeleteUserRetweet(context.Background(), a.CurrentUser.ID, tweetID)

	return checkError(err)
}
