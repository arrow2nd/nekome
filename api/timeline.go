package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchHomeTileline : ホームタイムラインを取得
func (a *API) FetchHomeTileline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.UserTweetReverseChronologicalTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.UserTweetReverseChronologicalTimeline(context.Background(), userID, opts)

	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if e := checkPartialError(res.Raw.Errors); e != nil {
		return nil, nil, e
	}

	return createTweetDictionarySlice(res.Raw), res.RateLimit, nil
}

// FetchUserTimeline : ユーザタイムラインを取得
func (a *API) FetchUserTimeline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.UserTweetTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.UserTweetTimeline(context.Background(), userID, opts)

	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if e := checkPartialError(res.Raw.Errors); e != nil {
		return nil, nil, e
	}

	return createTweetDictionarySlice(res.Raw), res.RateLimit, nil
}

// FetchUserMentionTimeline : ユーザのメンションタイムラインを取得
func (a *API) FetchUserMentionTimeline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.UserMentionTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := a.client.UserMentionTimeline(context.Background(), userID, opts)

	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if e := checkPartialError(res.Raw.Errors); e != nil {
		return nil, nil, e
	}

	return createTweetDictionarySlice(res.Raw), res.RateLimit, nil
}
