package app

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
)

type tweetActionType string
type userActionType string

const (
	like      tweetActionType = "Like"
	unlike    tweetActionType = "Unlike"
	retweet   tweetActionType = "Retweet"
	unretweet tweetActionType = "Unretweet"
	follow    userActionType  = "Follow"
	unfollow  userActionType  = "UnFollow"
	block     userActionType  = "Block"
	unblock   userActionType  = "Unblock"
	mute      userActionType  = "Mute"
	unmute    userActionType  = "Unmute"
)

// actionOnTweet : ツイートに対しての操作
func (t *tweets) actionOnTweet(a tweetActionType) {
	c := t.getSelectTweet()

	label := string(a)
	id := c.Tweet.ID
	summary := createTweetSummary(c)

	f := func() {
		var err error

		switch a {
		case like:
			err = shared.api.Like(id)
		case unlike:
			err = shared.api.UnLike(id)
		case retweet:
			err = shared.api.Retweet(id)
		case unretweet:
			err = shared.api.UnRetweet(id)
		}

		if err != nil {
			shared.SetErrorStatus(label, err.Error())
			return
		}

		shared.SetStatus(label+"ed", summary)
	}

	shared.ReqestPopupModal(&ModalOpt{
		fmt.Sprintf("Are you sure you want to %s this tweet?", strings.ToLower(label)),
		f,
	})
}

// actionOnUser : ユーザの操作
func (t *tweets) actionOnUser(a userActionType) {
	c := t.getSelectTweet()

	label := string(a)
	id := c.Author.ID
	summary := createUserSummary(c.Author)

	f := func() {
		var err error

		switch a {
		case follow:
			err = shared.api.Follow(id)
		case unfollow:
			err = shared.api.UnFollow(id)
		case block:
			err = shared.api.Block(id)
		case unblock:
			err = shared.api.UnBlock(id)
		case mute:
			err = shared.api.Mute(id)
		case unmute:
			err = shared.api.UnMute(id)
		}

		if err != nil {
			shared.SetErrorStatus(label, err.Error())
			return
		}

		shared.SetStatus(label+"ed", summary)
	}

	shared.ReqestPopupModal(&ModalOpt{
		fmt.Sprintf(`Are you sure you want to %s "%s"?`, strings.ToLower(label), summary),
		f,
	})
}

// openBrower : ブラウザで開く
func (t *tweets) openBrower() {
	c := t.getSelectTweet()

	if err := browser.OpenURL(createTweetURL(c)); err != nil {
		shared.SetErrorStatus("Open", err.Error())
		return
	}

	shared.SetStatus("Opened", createTweetSummary(c))
}

// copyLinkToClipBoard : リンクをクリップボードへコピー
func (t *tweets) copyLinkToClipBoard() {
	c := t.getSelectTweet()

	if err := clipboard.WriteAll(createTweetURL(c)); err != nil {
		shared.SetErrorStatus("Copy", err.Error())
		return
	}

	shared.SetStatus("Copied", createTweetSummary(c))
}
