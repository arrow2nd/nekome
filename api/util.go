package api

import "github.com/g8rswimmer/go-twitter/v2"

func createTweetDictionarySlice(raw *twitter.TweetRaw) []*twitter.TweetDictionary {
	contents := []*twitter.TweetDictionary{}
	dics := raw.TweetDictionaries()

	for _, tweet := range raw.Tweets {
		contents = append(contents, dics[tweet.ID])
	}

	return contents
}

// UserDictionary : 独自の twitter.UserDictionary 型
type UserDictionary struct {
	User        *twitter.UserObj
	PinnedTweet *twitter.TweetDictionary
}

func createUserDictionarySlice(raw *twitter.UserRaw) []*UserDictionary {
	users := []*UserDictionary{}
	dics := raw.UserDictionaries()

	for _, user := range raw.Users {
		var pinnedTweetDic *twitter.TweetDictionary = nil

		pinnedTweet := dics[user.ID].PinnedTweet

		// HACK: TweetObj を TweetDictionary に無理やり変換
		if pinnedTweet != nil {
			pinnedTweetDic = twitter.CreateTweetDictionary(*pinnedTweet, &twitter.TweetRawIncludes{
				Tweets: []*twitter.TweetObj{},
				Users:  []*twitter.UserObj{user},
				Places: []*twitter.PlaceObj{},
				Media:  []*twitter.MediaObj{},
				Polls:  []*twitter.PollObj{},
			})
		}

		users = append(users, &UserDictionary{
			User:        user,
			PinnedTweet: pinnedTweetDic,
		})
	}

	return users
}
