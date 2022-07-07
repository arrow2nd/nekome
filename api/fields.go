package api

import "github.com/g8rswimmer/go-twitter/v2"

var (
	// tweetFields : ツイートフィールド
	tweetFields = []twitter.TweetField{
		twitter.TweetFieldCreatedAt,
		twitter.TweetFieldAuthorID,
		twitter.TweetFieldPublicMetrics,
		twitter.TweetFieldEntities,
		twitter.TweetFieldSource,
		twitter.TweetFieldReferencedTweets,
	}

	// userFieldsForTL : タイムライン取得時のユーザフィールド
	userFieldsForTL = []twitter.UserField{
		twitter.UserFieldUserName,
		twitter.UserFieldName,
		twitter.UserFieldVerified,
		twitter.UserFieldProtected,
	}

	// userFieldsForUser : ユーザ詳細取得時のユーザフィールド
	userFieldsForUser = append(userFieldsForTL,
		twitter.UserFieldDescription,
		twitter.UserFieldLocation,
		twitter.UserFieldURL,
		twitter.UserFieldPublicMetrics,
	)

	// pollFields : 投票フィールド
	pollFields = []twitter.PollField{
		twitter.PollFieldVotingStatus,
		twitter.PollFieldEndDateTime,
	}

	// tweetExpansions : ツイートの拡張フィールド
	tweetExpansions = []twitter.Expansion{
		twitter.ExpansionAuthorID,
		twitter.ExpansionAttachmentsPollIDs,
		twitter.ExpansionReferencedTweetsID,
		twitter.ExpansionReferencedTweetsIDAuthorID,
	}
)
