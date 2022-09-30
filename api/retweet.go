package api

import (
	"context"
)

// Retweet : リツイート
func (a *API) Retweet(tweetId string) error {
	_, err := a.client.UserRetweet(context.Background(), a.CurrentUser.ID, tweetId)

	return checkError(err)
}

// UnRetweet : リツイートを解除
func (a *API) UnRetweet(tweetId string) error {
	_, err := a.client.DeleteUserRetweet(context.Background(), a.CurrentUser.ID, tweetId)

	return checkError(err)
}
