package app

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
)

// ツイートへの操作タイプ
const (
	tweetActionLike      string = "Like"
	tweetActionUnlike    string = "Unlike"
	tweetActionRetweet   string = "Retweet"
	tweetActionUnretweet string = "Unretweet"
	tweetActionDelete    string = "Delete"
)

// ユーザへの操作タイプ
const (
	userActionFollow   string = "Follow"
	userActionUnfollow string = "Unfollow"
	userActionBlock    string = "Block"
	userActionUnblock  string = "Unblock"
	userActionMute     string = "Mute"
	userActionUnmute   string = "Unmute"
)

// actionForTweet : ツイートへの操作
func (t *tweets) actionForTweet(a string) {
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
		case tweetActionLike:
			err = shared.api.Like(id)
		case tweetActionUnlike:
			err = shared.api.UnLike(id)
		case tweetActionRetweet:
			err = shared.api.Retweet(id)
		case tweetActionUnretweet:
			err = shared.api.UnRetweet(id)
		case tweetActionDelete:
			err = shared.api.DeleteTweet(id)
		default:
			shared.SetErrorStatus(label, fmt.Sprintf("unknown tweet action type: %s", a))
			return
		}

		if err != nil {
			shared.SetErrorStatus(label, err.Error())
			return
		}

		// TODO: RTの解除でもリストからツイートを削除するようにしたいが、
		//       RTしていなくてもAPIのリクエストが通ってしまうので未実装。

		// ツイート削除ならリストからツイートを削除
		if a == tweetActionDelete {
			t.DeleteTweet(id)
		}

		if !strings.HasSuffix(label, "e") {
			label += "e"
		}

		shared.SetStatus(label+"d", summary)
	}

	// 確認画面が不要ならそのまま実行
	if !shared.conf.Pref.Confirm[strings.ToLower(label)] {
		f()
		return
	}

	title := fmt.Sprintf(
		"Do you want to [%s]%s[-:-:-] this tweet?",
		shared.conf.Style.App.EmphasisText,
		strings.ToLower(label),
	)

	shared.ReqestPopupModal(&ModalOpt{title, "", f})
}

// actionForUser : ユーザへの操作
func (t *tweets) actionForUser(a string) {
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
		case userActionFollow:
			err = shared.api.Follow(id)
		case userActionUnfollow:
			err = shared.api.UnFollow(id)
		case userActionBlock:
			err = shared.api.Block(id)
		case userActionUnblock:
			err = shared.api.UnBlock(id)
		case userActionMute:
			err = shared.api.Mute(id)
		case userActionUnmute:
			err = shared.api.UnMute(id)
		default:
			shared.SetErrorStatus(label, fmt.Sprintf("unknown user action type: %s", a))
			return
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
	if !shared.conf.Pref.Confirm[strings.ToLower(label)] {
		f()
		return
	}

	title := fmt.Sprintf(
		`Do you want to [%s]%s[-:-:-] this user?`,
		shared.conf.Style.App.EmphasisText,
		strings.ToLower(label),
	)

	shared.ReqestPopupModal(&ModalOpt{title, summary, f})
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

// openUserLikes : ユーザのいいねリストを開く
func (t *tweets) openUserLikes() {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	cmd := fmt.Sprintf("likes %s", c.Author.UserName)
	shared.RequestExecCommand(cmd)
}

// insertQuoteCommand : Quoteコマンドを挿入
func (t *tweets) insertQuoteCommand() {
	c := t.getSelectTweet()
	if c == nil {
		return
	}

	cmd := fmt.Sprintf("tweet --quote %s ", c.Tweet.ID)
	shared.RequestInputCommand(cmd)
}

// insertReplyCommand : replyコマンドを挿入
func (t *tweets) insertReplyCommand() {
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

	url, err := createTweetUrl(c)
	if err != nil {
		shared.SetErrorStatus("Open", err.Error())
		return
	}

	if err := browser.OpenURL(url); err != nil {
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

	url, err := createTweetUrl(c)
	if err != nil {
		shared.SetErrorStatus("Copy", err.Error())
		return
	}

	if err := clipboard.WriteAll(url); err != nil {
		shared.SetErrorStatus("Copy", err.Error())
		return
	}

	shared.SetStatus("Copied", createTweetSummary(c))
}
