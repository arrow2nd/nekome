package api

import (
	"context"
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
)

// FetchHomeTileline ホームタイムラインを取得
func (a *API) FetchHomeTileline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserTweetReverseChronologicalTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFields,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := client.UserTweetReverseChronologicalTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("home timeline error: %v", err)
	}

	return createTweetDictionarySlice(res.Raw), nil
}

// FetchUserTimeline ユーザタイムラインを取得
func (a *API) FetchUserTimeline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserTweetTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFields,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := client.UserTweetTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("user timeline error: %v", err)
	}

	return createTweetDictionarySlice(res.Raw), nil
}

// FetchUserMentionTimeline ユーザのメンションタイムラインを取得
func (a *API) FetchUserMentionTimeline(userID, sinceID string, results int) ([]*twitter.TweetDictionary, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	client := a.newClient(a.CurrentUser.Token)

	opts := twitter.UserMentionTimelineOpts{
		TweetFields: tweetFields,
		PollFields:  pollFields,
		UserFields:  userFields,
		Expansions:  tweetExpansions,
		MaxResults:  results,
		SinceID:     sinceID,
	}

	res, err := client.UserMentionTimeline(context.Background(), userID, opts)
	if err != nil {
		return nil, fmt.Errorf("mention timeline error: %v", err)
	}

	return createTweetDictionarySlice(res.Raw), nil
}
