package app

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/pkg/browser"
)

// like : いいね
func (t *tweets) like() {
	c := t.getSelectTweet()
	f := func() {
		if err := shared.api.Like(c.Tweet.ID); err != nil {
			shared.SetErrorStatus("Like", err.Error())
			return
		}

		shared.SetStatus("Liked", createTweetSummary(c))
	}

	shared.ReqestPopupModal(&ModalOpt{
		"Are you sure you want to like this tweet?",
		f,
	})
}

// unLike : いいね解除
func (t *tweets) unLike() {
	c := t.getSelectTweet()
	f := func() {
		if err := shared.api.UnLike(c.Tweet.ID); err != nil {
			shared.SetErrorStatus("Unlike", c.Tweet.Text)
			return
		}

		shared.SetStatus("Unliked", createTweetSummary(c))
	}

	shared.ReqestPopupModal(&ModalOpt{
		"Are you sure you want to unlike this tweet?",
		f,
	})
}

// retweet : リツイート
func (t *tweets) retweet() {
	c := t.getSelectTweet()
	f := func() {
		if err := shared.api.Retweet(c.Tweet.ID); err != nil {
			shared.SetErrorStatus("Retweet", err.Error())
			return
		}

		shared.SetStatus("Retweeted", createTweetSummary(c))
	}

	shared.ReqestPopupModal(&ModalOpt{
		"Are you sure you want to retweet this tweet?",
		f,
	})
}

// unRetweet : リツイート解除
func (t *tweets) unRetweet() {
	c := t.getSelectTweet()
	f := func() {
		if err := shared.api.UnRetweet(c.Tweet.ID); err != nil {
			shared.SetErrorStatus("Unretweet", err.Error())
			return
		}

		shared.SetStatus("Unretweeted", createTweetSummary(c))
	}

	shared.ReqestPopupModal(&ModalOpt{
		"Are you sure you want to unretweet this tweet?",
		f,
	})
}

// follow : フォロー
func (t *tweets) follow() {
	c := t.getSelectTweet()
	s := createUserSummary(c.Author)

	f := func() {
		if err := shared.api.Follow(c.Author.ID); err != nil {
			shared.SetErrorStatus("Follow", err.Error())
			return
		}

		shared.SetStatus("Followed", s)
	}

	shared.ReqestPopupModal(&ModalOpt{
		fmt.Sprintf(`Are you sure you want to follow "%s"?`, s),
		f,
	})
}

// unfollow : フォロー解除
func (t *tweets) unfollow() {
	c := t.getSelectTweet()
	s := createUserSummary(c.Author)

	f := func() {
		if err := shared.api.Follow(c.Author.ID); err != nil {
			shared.SetErrorStatus("Unfollow", err.Error())
			return
		}

		shared.SetStatus("Unfollowed", s)
	}

	shared.ReqestPopupModal(&ModalOpt{
		fmt.Sprintf(`Are you sure you want to unfollow "%s"?`, s),
		f,
	})
}

// openBrower : ブラウザで表示
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
