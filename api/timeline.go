package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchHomeTileline ホームタイムラインを取得
func (a *API) FetchHomeTileline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, error) {
	opts := twitter.UserTweetReverseChronologicalTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFields,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.UserTweetReverseChronologicalTimeline(context.Background(), userID, opts)

	if e := checkError(err); e != nil {
		return nil, e
	}

	return createTweetDictionarySlice(res.Raw), nil
}

// FetchUserTimeline ユーザタイムラインを取得
func (a *API) FetchUserTimeline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, error) {
	opts := twitter.UserTweetTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFields,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.UserTweetTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("user timeline error: %v", err)
	}

	return createTweetDictionarySlice(res.Raw), nil
}

// FetchUserMentionTimeline ユーザのメンションタイムラインを取得
func (a *API) FetchUserMentionTimeline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, error) {
	opts := twitter.UserMentionTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFields,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.UserMentionTimeline(context.Background(), userID, opts)

	if e := checkError(err); e != nil {
		return nil, e
	}

	return createTweetDictionarySlice(res.Raw), nil
}
