package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// PostTweet : ツイートを投稿
func (a *API) PostTweet(text string) error {
	req := twitter.CreateTweetRequest{
		Text: text,
	}

	_, err := a.client.CreateTweet(context.Background(), req)

	return checkError(err)
}

// DeleteTweet : ツイートを削除
func (a *API) DeleteTweet(tweetId string) error {
	_, err := a.client.DeleteTweet(context.Background(), tweetId)

	return checkError(err)
}
