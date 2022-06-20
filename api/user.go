package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchUser : UserNameからユーザ情報を取得
func (a *API) FetchUser(userNames []string) ([]*UserDictionary, error) {
	opts := twitter.UserLookupOpts{
		TweetFields: tweetFields,
		UserFields:  userFieldsForUser,
		Expansions: []twitter.Expansion{
			twitter.ExpansionPinnedTweetID,
		},
	}

	res, err := a.client.UserNameLookup(context.Background(), userNames, opts)
	if err != nil {
		return nil, fmt.Errorf("username lookup error: %v", err)
	}

	return createUserDictionarySlice(res.Raw), nil
}
