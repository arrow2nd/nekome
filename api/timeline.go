package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchUserTimeline ユーザタイムラインを取得
func (a *API) FetchUserTimeline(userID string, results int) ([]*twitter.TweetObj, error) {
	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserTweetTimelineOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	result, err := client.UserTweetTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("user timeline error: %v", err)
	}

	return result.Raw.Tweets, nil
}

// FetchUserMentionTimeline ユーザのメンションタイムラインを取得
func (a *API) FetchUserMentionTimeline(userID string, results int) ([]*twitter.TweetObj, error) {
	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserMentionTimelineOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	result, err := client.UserMentionTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("mention timeline error: %v", err)
	}

	return result.Raw.Tweets, nil
}
