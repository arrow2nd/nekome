package api

import (
	"errors"
	"strconv"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/manifoldco/promptui"
)

var (
	consumerToken  = ""
	consumerSecret = ""
)

// inputPinCode : PINの入力を受付
func inputPinCode() (string, error) {
	prompt := promptui.Prompt{
		Label: "PIN",
		Validate: func(s string) error {
			if _, err := strconv.Atoi(s); err != nil {
				return errors.New("please enter a number")
			}
			return nil
		},
	}

	return prompt.Run()
}

// getConsumerToken : コンシューマトークンを取得
func getConsumerToken(ct *oauth1.Token) (*oauth1.Token, error) {
	consumer := &oauth1.Token{
		Token:       consumerToken,
		TokenSecret: consumerSecret,
	}

	if ct.Token != "" {
		consumer.Token = ct.Token
	}

	if ct.TokenSecret != "" {
		consumer.TokenSecret = ct.TokenSecret
	}

	if consumer.Token == "" || consumer.TokenSecret == "" {
		return nil, errors.New("no consumer key, please run 'nekome edit' and set it to .cred.toml")
	}

	return consumer, nil
}

// createTweetSlice : TweetDictionary のスライスを作成
func createTweetSlice(raw *twitter.TweetRaw) (bool, []*twitter.TweetDictionary) {
	// データがあるか
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

// UserDictionary : twitter.UserDictionary の独自実装
type UserDictionary struct {
	User        *twitter.UserObj
	PinnedTweet *twitter.TweetDictionary
}

// createUserSlice : UserDictionary のスライスを作成
func createUserSlice(raw *twitter.UserRaw) (bool, []*UserDictionary) {
	// データがあるか
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
