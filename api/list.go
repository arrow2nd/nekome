package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchOwnedListIDs : ユーザが所有するリストの情報を取得
func (a *API) FetchOwnedLists(userID string) ([]*twitter.ListObj, error) {
	opts := twitter.UserListLookupOpts{}

	result, err := a.client.UserListLookup(context.Background(), userID, opts)
	if e := checkError(err); e != nil {
		return nil, e
	}

	return result.Raw.Lists, nil
}

// FetchListTweets : リストのツイートを取得
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

	return createTweetDictionarySlice(res.Raw), res.RateLimit, nil
}
