package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchOwnedLists : ユーザが所有するリストの情報を取得
func (a *API) FetchOwnedLists(userID string) ([]*twitter.ListObj, error) {
	opts := twitter.UserListLookupOpts{}

	res, err := a.client.UserListLookup(context.Background(), userID, opts)
	if e := checkError(err); e != nil {
		return nil, e
	}

	if res.Raw == nil {
		return []*twitter.ListObj{}, nil
	}

	if e := checkPartialError(res.Raw.Errors); len(res.Raw.Lists) == 0 && e != nil {
		return nil, e
	}

	return res.Raw.Lists, nil
}

// FetchListTweets : リスト内のツイートを取得
func (a *API) FetchListTweets(listID string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.ListTweetLookupOpts{
		TweetFields: tweetFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
	}

	res, err := a.client.ListTweetLookup(context.Background(), listID, opts)
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
