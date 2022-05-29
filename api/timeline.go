package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// UserTimeline ユーザタイムラインを取得
func (a *API) UserTimeline(userID string, results int) ([]*twitter.TweetObj, error) {
	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserTweetTimelineOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	timeline, err := client.UserTweetTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("user timeline error: %v", err)
	}

	return timeline.Raw.Tweets, nil
}

// UserMentionTimeline ユーザのメンションタイムラインを取得
func (a *API) UserMentionTimeline(userID string, results int) ([]*twitter.TweetObj, error) {
	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserMentionTimelineOpts{
		TweetFields: tweetFields,
		UserFields:  userFields,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	timeline, err := client.UserMentionTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("mention timeline error: %v", err)
	}

	return timeline.Raw.Tweets, nil
}
