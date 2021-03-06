package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// PostTweet : ツイートを投稿
func (a *API) PostTweet(text, quoteId, replyId string, mediaIDs []string) error {
	req := twitter.CreateTweetRequest{
		Text: text,
	}

	if quoteId != "" {
		req.QuoteTweetID = quoteId
	}

	if replyId != "" {
		req.Reply = &twitter.CreateTweetReply{
			InReplyToTweetID: replyId,
		}
	}

	if len(mediaIDs) > 0 {
		req.Media = &twitter.CreateTweetMedia{
			IDs: mediaIDs,
		}
	}

	_, err := a.client.CreateTweet(context.Background(), req)

	return checkError(err)
}

// DeleteTweet : ツイートを削除
func (a *API) DeleteTweet(tweetId string) error {
	_, err := a.client.DeleteTweet(context.Background(), tweetId)

	return checkError(err)
}
