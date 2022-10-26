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
func createTweetSlice(raw *twitter.TweetRaw) ([]*twitter.TweetDictionary, bool) {
	// データがあるか
	if len(raw.Tweets) == 0 || raw.Tweets[0] == nil {
		return nil, false
	}

	contents := []*twitter.TweetDictionary{}
	dics := raw.TweetDictionaries()

	for _, tweet := range raw.Tweets {
		contents = append(contents, dics[tweet.ID])
	}

	return contents, true
}

// UserDictionary : twitter.UserDictionary の独自実装
type UserDictionary struct {
	User        *twitter.UserObj
	PinnedTweet *twitter.TweetDictionary
}

// createUserSlice : UserDictionary のスライスを作成
func createUserSlice(raw *twitter.UserRaw, pinnedTweetRaw *twitter.TweetRaw) ([]*UserDictionary, bool) {
	// データがあるか
	if len(raw.Users) == 0 || raw.Users[0] == nil {
		return nil, false
	}

	pinnedTweets := map[string]*twitter.TweetDictionary{}

	// NOTE: nilでないかつ、部分エラーが発生していないならピン止めツイートがあるとみなす
	existPinnedTweet := pinnedTweetRaw != nil && checkPartialError(pinnedTweetRaw.Errors) == nil
	if existPinnedTweet {
		pinnedTweets = pinnedTweetRaw.TweetDictionaries()
	}

	users := []*UserDictionary{}
	for _, u := range raw.Users {
		dictionary := &UserDictionary{
			User:        u,
			PinnedTweet: nil,
		}

		if u.PinnedTweetID != "" {
			dictionary.PinnedTweet = pinnedTweets[u.PinnedTweetID]
		}

		users = append(users, dictionary)
	}

	return users, true
}
