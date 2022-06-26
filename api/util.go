package api

import (
	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

var (
	consumerToken  = ""
	consumerSecret = ""
)

// getConsumerToken : クライアントトークンを取得
func getConsumerToken(ct *oauth1.Token) *oauth1.Token {
	token := consumerToken
	secret := consumerSecret

	if ct.Token != "" {
		token = ct.Token
	}

	if ct.TokenSecret != "" {
		secret = ct.TokenSecret
	}

	return &oauth1.Token{
		Token:       token,
		TokenSecret: secret,
	}
}

func createTweetDictionarySlice(raw *twitter.TweetRaw) (bool, []*twitter.TweetDictionary) {
	// データがあるかチェック
	if len(raw.Tweets) == 0 || raw.Tweets[0] == nil {
		return false, nil
	}

	contents := []*twitter.TweetDictionary{}
	dics := raw.TweetDictionaries()

	for _, tweet := range raw.Tweets {
		contents = append(contents, dics[tweet.ID])
	}

	return true, contents
}

// UserDictionary : 独自の twitter.UserDictionary 型
type UserDictionary struct {
	User        *twitter.UserObj
	PinnedTweet *twitter.TweetDictionary
}

func createUserDictionarySlice(raw *twitter.UserRaw) (bool, []*UserDictionary) {
	// データがあるかチェック
	if len(raw.Users) == 0 || raw.Users[0] == nil {
		return false, nil
	}

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

	return true, users
}
