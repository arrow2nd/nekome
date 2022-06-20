package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// PostTweet : ツイートを投稿
func (a *API) PostTweet(text string) error {
	req := twitter.CreateTweetRequest{
		Text: text,
	}

	if _, err := a.client.CreateTweet(context.Background(), req); err != nil {
		return fmt.Errorf("post tweet error: %v", err)
	}

	return nil
}

// DeleteTweet : ツイートを削除
func (a *API) DeleteTweet(tweetId string) error {
	if _, err := a.client.DeleteTweet(context.Background(), tweetId); err != nil {
		return fmt.Errorf("delete tweet error: %v", err)
	}

	return nil
}
