package api

import (
	"context"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchHomeTileline : ホームタイムラインを取得
func (a *API) FetchHomeTileline(userId, sinceId string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.UserTweetReverseChronologicalTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceId,
	}

	res, err := a.client.UserTweetReverseChronologicalTimeline(context.Background(), userId, opts)
	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if res.Raw == nil {
		return []*twitter.TweetDictionary{}, res.RateLimit, nil
	}

	ok, tweets := createTweetSlice(res.Raw)
	if e := checkPartialError(res.Raw.Errors); !ok && e != nil {
		return nil, nil, e
	}

	return tweets, res.RateLimit, nil
}

// FetchUserTimeline : ユーザタイムラインを取得
func (a *API) FetchUserTimeline(userId, sinceId string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.UserTweetTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceId,
	}

	res, err := a.client.UserTweetTimeline(context.Background(), userId, opts)
	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if res.Raw == nil {
		return []*twitter.TweetDictionary{}, res.RateLimit, nil
	}

	ok, tweets := createTweetSlice(res.Raw)
	if e := checkPartialError(res.Raw.Errors); !ok && e != nil {
		return nil, nil, e
	}

	return tweets, res.RateLimit, nil
}

// FetchUserMentionTimeline : ユーザのメンションタイムラインを取得
func (a *API) FetchUserMentionTimeline(userId, sinceId string, results int) ([]*twitter.TweetDictionary, *twitter.RateLimit, error) {
	opts := twitter.UserMentionTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFieldsForTL,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceId,
	}

	res, err := a.client.UserMentionTimeline(context.Background(), userId, opts)

	if e := checkError(err); e != nil {
		return nil, nil, e
	}

	if res.Raw == nil {
		return []*twitter.TweetDictionary{}, res.RateLimit, nil
	}

	ok, tweets := createTweetSlice(res.Raw)
	if e := checkPartialError(res.Raw.Errors); !ok && e != nil {
		return nil, nil, e
	}

	return tweets, res.RateLimit, nil
}
