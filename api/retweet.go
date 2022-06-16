package api

import (
	"context"
	"fmt"
)

// Retweet リツイート
func (a *API) Retweet(tweetID string) error {
	if _, err := a.client.UserRetweet(context.Background(), a.CurrentUser.ID, tweetID); err != nil {
		return fmt.Errorf("retweet error: %v", err)
	}

	return nil
}

// UnRetweet リツイートを解除
func (a *API) UnRetweet(tweetID string) error {
	if _, err := a.client.DeleteUserRetweet(context.Background(), a.CurrentUser.ID, tweetID); err != nil {
		return fmt.Errorf("unretweet error: %v", err)
	}
	return nil
}
