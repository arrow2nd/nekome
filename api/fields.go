package api

import "github.com/g8rswimmer/go-twitter/v2"

var (
	tweetFields = []twitter.TweetField{
		twitter.TweetFieldCreatedAt,
		twitter.TweetFieldAuthorID,
		twitter.TweetFieldPublicMetrics,
		twitter.TweetFieldEntities,
	}

	userFields = []twitter.UserField{
		twitter.UserFieldUserName,
		twitter.UserFieldName,
	}
)
