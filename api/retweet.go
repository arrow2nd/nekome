package api

import (
	"context"
	"fmt"
)

// Retweet リツイート
func (a *API) Retweet(tweetId string) error {
	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.UserRetweet(context.Background(), a.CurrentUser.ID, tweetId); err != nil {
		return fmt.Errorf("retweet error: %v", err)
	}

	return nil
}
