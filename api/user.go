package api

import (
	"context"

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
	if e := checkError(err); e != nil {
		return nil, e
	}

	if res.Raw == nil {
		return []*UserDictionary{}, nil
	}

	ok, users := createUserSlice(res.Raw)
	if e := checkPartialError(res.Raw.Errors); !ok && e != nil {
		return nil, e
	}

	return users, nil
}
