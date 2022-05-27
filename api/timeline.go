package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchUserTimeline ユーザタイムラインを取得
func (a *API) FetchUserTimeline(userID string) ([]*twitter.TweetObj, error) {
	client, err := a.newClient(a.token)
	if err != nil {
		return nil, err
	}

	opts := twitter.UserTweetTimelineOpts{
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldAuthorID,
			twitter.TweetFieldPublicMetrics,
			twitter.TweetFieldEntities,
		},
		UserFields: []twitter.UserField{
			twitter.UserFieldUserName,
			twitter.UserFieldName,
		},
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: 25,
	}

	timeline, err := client.UserTweetTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, err
	}

	return timeline.Raw.Tweets, nil
}
