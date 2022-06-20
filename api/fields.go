package api

import "github.com/g8rswimmer/go-twitter/v2"

var (
	tweetFields = []twitter.TweetField{
		twitter.TweetFieldCreatedAt,
		twitter.TweetFieldAuthorID,
		twitter.TweetFieldPublicMetrics,
		twitter.TweetFieldEntities,
		twitter.TweetFieldSource,
		twitter.TweetFieldReferencedTweets,
	}

	userFieldsForTL = []twitter.UserField{
		twitter.UserFieldUserName,
		twitter.UserFieldName,
		twitter.UserFieldVerified,
		twitter.UserFieldProtected,
	}

	userFieldsForUser = append(userFieldsForTL,
		twitter.UserFieldDescription,
		twitter.UserFieldLocation,
		twitter.UserFieldURL,
		twitter.UserFieldPublicMetrics,
	)

	pollFields = []twitter.PollField{
		twitter.PollFieldVotingStatus,
		twitter.PollFieldEndDateTime,
	}

	tweetExpansions = []twitter.Expansion{
		twitter.ExpansionAuthorID,
		twitter.ExpansionAttachmentsPollIDs,
		twitter.ExpansionReferencedTweetsID,
		twitter.ExpansionReferencedTweetsIDAuthorID,
	}
)
