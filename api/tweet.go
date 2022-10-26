package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// PostTweet : ツイートを投稿
func (a *API) PostTweet(text, quoteID, replyID string, mediaIDs []string) error {
	req := twitter.CreateTweetRequest{
		Text: text,
	}

	if quoteID != "" {
		req.QuoteTweetID = quoteID
	}

	if replyID != "" {
		req.Reply = &twitter.CreateTweetReply{
			InReplyToTweetID: replyID,
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
func (a *API) DeleteTweet(tweetID string) error {
	_, err := a.client.DeleteTweet(context.Background(), tweetID)

	return checkError(err)
}
