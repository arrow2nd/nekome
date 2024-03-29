package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchLikedTweets : ユーザのいいねしたツイートを取得
func (a *API) FetchLikedTweets(userID string, maxResults int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.UserLikesLookupOpts{
		TweetFields: tweetFields,
		UserFields:  userFieldsForTL,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: maxResults,
	}
	res, err := a.client.UserLikesLookup(context.Background(), userID, opts)
	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if res.Raw == nil {
		return []*twitter.TweetDictionary{}, res.RateLimit, nil
	}

	tweets, ok := createTweetSlice(res.Raw)
	if e := checkPartialError(res.Raw.Errors); !ok && e != nil {
		return nil, nil, e
	}

	return tweets, res.RateLimit, nil
}

// Like : いいね
func (a *API) Like(tweetID string) error {
	_, err := a.client.UserLikes(context.Background(), a.CurrentUser.ID, tweetID)

	return checkError(err)
}

// UnLike : いいねを解除
func (a *API) UnLike(tweetID string) error {
	_, err := a.client.DeleteUserLikes(context.Background(), a.CurrentUser.ID, tweetID)

	return checkError(err)
}
