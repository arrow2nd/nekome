package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// PostTweet ツイートを投稿
func (a *API) PostTweet(text string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	req := twitter.CreateTweetRequest{
		Text: text,
	}

	if _, err := client.CreateTweet(context.Background(), req); err != nil {
		return fmt.Errorf("post tweet error: %v", err)
	}

	return nil
}

// DeleteTweet ツイートを削除
func (a *API) DeleteTweet(tweetId string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	if _, err := client.DeleteTweet(context.Background(), tweetId); err != nil {
		return fmt.Errorf("delete tweet error: %v", err)
	}

	return nil
}
