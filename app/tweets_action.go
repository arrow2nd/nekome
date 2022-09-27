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
	tweetLike      tweetActionType = "Like"
	tweetUnlike    tweetActionType = "Unlike"
	tweetRetweet   tweetActionType = "Retweet"
	tweetUnretweet tweetActionType = "Unretweet"
	tweetDelete    tweetActionType = "Delete"
	userFollow     userActionType  = "Follow"
	userUnfollow   userActionType  = "Unfollow"
	userBlock      userActionType  = "Block"
	userUnblock    userActionType  = "Unblock"
	userMute       userActionType  = "Mute"
	userUnmute     userActionType  = "Unmute"
)

// actionForTweet : ツイートに対しての操作
func (t *tweets) actionForTweet(a tweetActionType) {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	label := string(a)
	id := c.Tweet.ID
	summary := createTweetSummary(c)

	f := func() {
		var err error

		switch a {
		case tweetLike:
			err = shared.api.Like(id)
		case tweetUnlike:
			err = shared.api.UnLike(id)
		case tweetRetweet:
			err = shared.api.Retweet(id)
		case tweetUnretweet:
			err = shared.api.UnRetweet(id)
		case tweetDelete:
			err = shared.api.DeleteTweet(id)
		}

		if err != nil {
			shared.SetErrorStatus(label, err.Error())
			return
		}

		if !strings.HasSuffix(label, "e") {
			label += "e"
		}

		shared.SetStatus(label+"d", summary)
	}

	// 確認画面が不要ならそのまま実行
	if !shared.conf.Settings.Confirm[strings.ToLower(label)] {
		f()
		return
	}

	style := shared.conf.Style.App.EmphasisText

	shared.ReqestPopupModal(&ModalOpt{
		fmt.Sprintf("Do you want to [%s]%s[-:-:-] this tweet?", style, strings.ToLower(label)),
		"",
		f,
	})
}

// actionForUser : ユーザに対しての操作
func (t *tweets) actionForUser(a userActionType) {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	label := string(a)
	id := c.Author.ID
	summary := createUserSummary(c.Author)

	f := func() {
		var err error

		switch a {
		case userFollow:
			err = shared.api.Follow(id)
		case userUnfollow:
			err = shared.api.UnFollow(id)
		case userBlock:
			err = shared.api.Block(id)
		case userUnblock:
			err = shared.api.UnBlock(id)
		case userMute:
			err = shared.api.Mute(id)
		case userUnmute:
			err = shared.api.UnMute(id)
		}

		if err != nil {
			shared.SetErrorStatus(label, err.Error())
			return
		}

		if !strings.HasSuffix(label, "e") {
			label += "e"
		}

		shared.SetStatus(label+"d", summary)
	}

	// 確認画面が不要ならそのまま実行
	if !shared.conf.Settings.Confirm[strings.ToLower(label)] {
		f()
		return
	}

	style := shared.conf.Style.App.EmphasisText

	shared.ReqestPopupModal(&ModalOpt{
		fmt.Sprintf(`Do you want to [%s]%s[-:-:-] this user?`, style, strings.ToLower(label)),
		summary,
		f,
	})
}

// openUserPage : ユーザページを開く
func (t *tweets) openUserPage() {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	cmd := fmt.Sprintf("user %s", c.Author.UserName)
	shared.RequestExecCommand(cmd)
}

// postQuoteTweet : QTコマンドを挿入
func (t *tweets) postQuoteTweet() {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	cmd := fmt.Sprintf("tweet --quote %s ", c.Tweet.ID)
	shared.RequestInputCommand(cmd)
}

// postReply : リプライコマンドを挿入
func (t *tweets) postReply() {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	cmd := fmt.Sprintf("tweet --reply %s ", c.Tweet.ID)
	shared.RequestInputCommand(cmd)
}

// openBrower : ブラウザで開く
func (t *tweets) openBrower() {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	if err := browser.OpenURL(createTweetURL(c)); err != nil {
		shared.SetErrorStatus("Open", err.Error())
		return
	}

	shared.SetStatus("Opened", createTweetSummary(c))
}

// copyLinkToClipBoard : リンクをクリップボードへコピー
func (t *tweets) copyLinkToClipBoard() {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	if err := clipboard.WriteAll(createTweetURL(c)); err != nil {
		shared.SetErrorStatus("Copy", err.Error())
		return
	}

	shared.SetStatus("Copied", createTweetSummary(c))
}
