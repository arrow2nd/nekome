package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

var (
	tweetFiled = []twitter.TweetField{
		twitter.TweetFieldCreatedAt,
		twitter.TweetFieldAuthorID,
		twitter.TweetFieldPublicMetrics,
		twitter.TweetFieldEntities,
	}
	userFiled = []twitter.UserField{
		twitter.UserFieldUserName,
		twitter.UserFieldName,
	}
)

// UserTimeline ユーザタイムラインを取得
func (a *API) UserTimeline(userID string, results int) ([]*twitter.TweetObj, error) {
	client, err := a.newClient(a.CurrentUser.Token)
	if err != nil {
		return nil, err
	}

	opts := twitter.UserTweetTimelineOpts{
		TweetFields: tweetFiled,
		UserFields:  userFiled,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	timeline, err := client.UserTweetTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, err
	}

	return timeline.Raw.Tweets, nil
}

// UserMentionTimeline ユーザのメンションタイムラインを取得
func (a *API) UserMentionTimeline(userID string, results int) ([]*twitter.TweetObj, error) {
	client, err := a.newClient(a.CurrentUser.Token)
	if err != nil {
		return nil, err
	}

	opts := twitter.UserMentionTimelineOpts{
		TweetFields: tweetFiled,
		UserFields:  userFiled,
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
		},
		MaxResults: results,
	}

	timeline, err := client.UserMentionTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, err
	}

	return timeline.Raw.Tweets, nil
}
