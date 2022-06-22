package ui

import (
	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
)

// like : いいね
func (t *tweets) like() {
	c := t.getSelectTweet()

	err := shared.api.Like(c.Tweet.ID)
	if err != nil {
		shared.SetErrorStatus("Like", err.Error())
		return
	}

	shared.SetStatus("Liked", createTweetSummary(c))
}

// unLike : いいね解除
func (t *tweets) unLike() {
	c := t.getSelectTweet()

	err := shared.api.UnLike(c.Tweet.ID)
	if err != nil {
		shared.SetErrorStatus("UnLike", c.Tweet.Text)
		return
	}

	shared.SetStatus("UnLiked", createTweetSummary(c))
}

// retweet : リツイート
func (t *tweets) retweet() {
	c := t.getSelectTweet()

	err := shared.api.Retweet(c.Tweet.ID)
	if err != nil {
		shared.SetErrorStatus("Retweet", err.Error())
		return
	}

	shared.SetStatus("Retweeted", createTweetSummary(c))
}

// unRetweet : リツイート解除
func (t *tweets) unRetweet() {
	c := t.getSelectTweet()

	err := shared.api.UnRetweet(c.Tweet.ID)
	if err != nil {
		shared.SetErrorStatus("UnRetweet", err.Error())
		return
	}

	shared.SetStatus("UnRetweeted", createTweetSummary(c))
}

// openBrower : ブラウザで表示
func (t *tweets) openBrower() {
	shared.SetStatus("Open", "Wait...")

	c := t.getSelectTweet()

	browser.OpenURL(createTweetURL(c))

	shared.SetStatus("Opened", createTweetSummary(c))
}

// copyLinkToClipBoard : リンクをクリップボードへコピー
func (t *tweets) copyLinkToClipBoard() {
	c := t.getSelectTweet()

	clipboard.WriteAll(createTweetURL(c))

	shared.SetStatus("Copied", createTweetSummary(c))
}
