package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchUser : UserNameからユーザ情報を取得
func (a *API) FetchUser(userNames []string) ([]*UserDictionary, error) {
	opts := twitter.UserLookupOpts{
		UserFields: userFieldsForUser,
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

	pinnedTweetRew, err := a.fetchPinnedTweets(res.Raw)
	if e := checkError(err); e != nil {
		return nil, e
	}

	users, ok := createUserSlice(res.Raw, pinnedTweetRew)
	if e := checkPartialError(res.Raw.Errors); !ok && e != nil {
		return nil, e
	}

	return users, nil
}

func (a *API) fetchPinnedTweets(raw *twitter.UserRaw) (*twitter.TweetRaw, error) {
	tweetIDs := []string{}
	for _, user := range raw.Users {
		if user != nil && user.PinnedTweetID != "" {
			tweetIDs = append(tweetIDs, user.PinnedTweetID)
		}
	}

	// ピン止めツイートが無い
	if len(tweetIDs) == 0 {
		return nil, nil
	}

	opts := twitter.TweetLookupOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
	}

	res, err := a.client.TweetLookup(context.Background(), tweetIDs, opts)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return res.Raw, nil
}
