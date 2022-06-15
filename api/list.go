package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchOwnedListIDs ユーザが所有するリストの情報を取得
func (a *API) FetchOwnedLists(userID string) ([]*twitter.ListObj, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserListLookupOpts{}

	result, err := client.UserListLookup(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("user owned list lookup error: %v", err)
	}

	return result.Raw.Lists, nil
}

// FetchListTweets リストのツイートを取得
func (a *API) FetchListTweets(listID string, results int) ([]*twitter.TweetDictionary, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.ListTweetLookupOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	res, err := client.ListTweetLookup(context.Background(), listID, opts)
	if err != nil {
		return nil, fmt.Errorf("list tweet lookup error: %v", err)
	}

	return createTweetDictionarySlice(res.Raw), nil
}
