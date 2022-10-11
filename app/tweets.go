package app

import (
	"fmt"
	"sync"

	"github.com/arrow2nd/nekome/v2/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rivo/tview"
)

// cursorMove : カーソルの移動量
const (
	cursorMoveUp   int = -1
	cursorMoveDown int = 1
)

// tweets : ツイートの表示管理
type tweets struct {
	view      *tview.TextView
	pinned    *twitter.TweetDictionary
	contents  []*twitter.TweetDictionary
	rateLimit *twitter.RateLimit
	mu        sync.Mutex
}

func newTweets() (*tweets, error) {
	t := &tweets{
		view:      tview.NewTextView(),
		pinned:    nil,
		rateLimit: nil,
		contents:  []*twitter.TweetDictionary{},
	}

	t.view.
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true)

	t.view.SetHighlightedFunc(func(_, _, _ []string) {
		t.view.ScrollToHighlight()
	})

	if err := t.setKeybindings(); err != nil {
		return nil, err
	}

	return t, nil
}

// setKeybindings : キーバインドを設定
func (t *tweets) setKeybindings() error {
	handlers := map[string]func(){
		config.ActionScrollUp: func() {
			r, c := t.view.GetScrollOffset()
			t.view.ScrollTo(r+1, c)
		},
		config.ActionScrollDown: func() {
			r, c := t.view.GetScrollOffset()
			t.view.ScrollTo(r-1, c)
		},
		config.ActionCursorUp: func() {
			t.moveCursor(cursorMoveUp)
		},
		config.ActionCursorDown: func() {
			t.moveCursor(cursorMoveDown)
		},
		config.ActionCursorTop: func() {
			t.scrollToTweet(0)
		},
		config.ActionCursorBottom: func() {
			lastIndex := t.GetTweetsCount() - 1
			t.scrollToTweet(lastIndex)
		},
		config.ActionTweetLike: func() {
			t.actionForTweet(tweetActionLike)
		},
		config.ActionTweetUnlike: func() {
			t.actionForTweet(tweetActionUnlike)
		},
		config.ActionTweetRetweet: func() {
			t.actionForTweet(tweetActionRetweet)
		},
		config.ActionTweetUnretweet: func() {
			t.actionForTweet(tweetActionUnretweet)
		},
		config.ActionTweetDelete: func() {
			t.actionForTweet(tweetActionDelete)
		},
		config.ActionUserFollow: func() {
			t.actionForUser(userActionFollow)
		},
		config.ActionUserUnfollow: func() {
			t.actionForUser(userActionUnfollow)
		},
		config.ActionUserBlock: func() {
			t.actionForUser(userActionBlock)
		},
		config.ActionUserUnblock: func() {
			t.actionForUser(userActionUnblock)
		},
		config.ActionUserMute: func() {
			t.actionForUser(userActionMute)
		},
		config.ActionUserUnmute: func() {
			t.actionForUser(userActionUnmute)
		},
		config.ActionOpenUserPage: func() {
			t.openUserPage()
		},
		config.ActionOpenUserLikes: func() {
			t.openUserLikes()
		},
		config.ActionTweet: func() {
			shared.RequestExecCommand("tweet")
		},
		config.ActionQuote: func() {
			t.insertQuoteCommand()
		},
		config.ActionReply: func() {
			t.insertReplyCommand()
		},
		config.ActionOpenBrowser: func() {
			t.openBrower()
		},
		config.ActionCopyUrl: func() {
			t.copyLinkToClipBoard()
		},
	}

	c, err := shared.conf.Pref.Keybindings.TweetView.MappingEventHandler(handlers)
	if err != nil {
		return err
	}

	t.view.SetInputCapture(c.Capture)

	return nil
}

// GetSinceId : 一番新しいツイートIDを取得
func (t *tweets) GetSinceId() string {
	if len(t.contents) == 0 {
		return ""
	}

	return t.contents[0].Tweet.ID
}

// GetTweetsCount : 表示中のツイート数を取得
func (t *tweets) GetTweetsCount() int {
	c := len(t.contents)

	if t.pinned != nil {
		c++
	}

	return c
}

// getSelectTweet : 選択中のツイートを取得
func (t *tweets) getSelectTweet() *twitter.TweetDictionary {
	id := getHighlightId(t.view.GetHighlights())
	if id == -1 {
		return nil
	}

	c := &twitter.TweetDictionary{}

	if t.pinned == nil {
		// ピン留めツイートが無い
		c = t.contents[id]
	} else if id == 0 {
		// ピン留めツイートが選択されている
		c = t.pinned
	} else {
		// ピン留めツイート以外が選択されている
		c = t.contents[id-1]
	}

	// RTなら参照先に置き換え
	for _, rc := range c.ReferencedTweets {
		if rc.Reference.Type == "retweeted" {
			c = rc.TweetDictionary
		}
	}

	return c
}

// getCurrentCursorPos : 現在のカーソル位置を取得
func (t *tweets) getCurrentCursorPos() int {
	pos := getHighlightId(t.view.GetHighlights())

	if pos == -1 {
		pos = 0
	}

	return pos
}

// registerTweets : ツイートを登録
func (t *tweets) registerTweets(tweets []*twitter.TweetDictionary) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	size := len(t.contents)
	addSize := len(tweets)
	allSize := size + addSize
	maxSize := shared.conf.Pref.Feature.AccmulateTweetsLimit

	// 最大蓄積数を超えていたら古いものから削除
	if allSize > maxSize {
		size -= allSize - maxSize
	}

	t.contents = append(tweets, t.contents[:size]...)

	return addSize
}

// RegisterPinnedTweet : ピン留めツイートを登録
func (t *tweets) RegisterPinnedTweet(tweet *twitter.TweetDictionary) {
	t.pinned = tweet
}

// DeleteTweet : 該当ツイートをリストから削除
func (t *tweets) DeleteTweet(tweetId string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	i, ok := find(t.contents, func(c *twitter.TweetDictionary) bool {
		// リツイート先のIDを参照
		for _, rc := range c.ReferencedTweets {
			if rc.Reference.Type == "retweeted" {
				return rc.TweetDictionary.Tweet.ID == tweetId
			}
		}

		return c.Tweet.ID == tweetId
	})

	if !ok {
		return
	}

	// i番目の要素を削除
	t.contents = t.contents[:i+copy(t.contents[i:], t.contents[i+1:])]

	// 再描画して反映
	t.draw(t.getCurrentCursorPos())
}

// UpdateRateLimit : レート制限を更新
func (t *tweets) UpdateRateLimit(r *twitter.RateLimit) {
	if r != nil {
		t.rateLimit = r
	}
}

// Update : ツイートを更新
func (t *tweets) Update(tweets []*twitter.TweetDictionary) {
	addedTweetsCount := t.registerTweets(tweets)
	cursorPos := t.getCurrentCursorPos()

	// 先頭以外のツイートを選択中の場合、更新後もそのツイートを選択したままにする
	// NOTE: "先頭以外" なのは、ストリームモードで放置した時にカーソルが段々下に下がってしまうのを防ぐため
	if cursorPos != 0 {
		cursorPos += addedTweetsCount
	}

	t.draw(cursorPos)
}

// draw : 描画（表示幅はターミナルのウィンドウ幅に依存）
func (t *tweets) draw(cursorPos int) {
	icon := shared.conf.Pref.Icon
	appearance := shared.conf.Pref.Appearance
	width := getWindowWidth()

	t.view.
		SetTextAlign(tview.AlignLeft).
		Clear()

	// 表示するツイートが無いなら描画を中断
	if t.GetTweetsCount() == 0 {
		t.DrawMessage(shared.conf.Pref.Text.NoTweets)
		return
	}

	contents := t.contents

	// ピン留めツイートがある場合、先頭に追加
	if t.pinned != nil {
		contents = append([]*twitter.TweetDictionary{t.pinned}, t.contents...)
	}

	for i, content := range contents {
		var quotedTweet *twitter.TweetDictionary = nil
		annotation := ""

		// 参照ツイートを確認
		for _, rc := range content.ReferencedTweets {
			switch rc.Reference.Type {
			case "retweeted":
				annotation += createAnnotation("RT by", content.Author)
				content = rc.TweetDictionary
			case "replied_to":
				annotation += createAnnotation("Reply to", rc.TweetDictionary.Author)
			case "quoted":
				quotedTweet = rc.TweetDictionary
			}
		}

		// ピン留めツイート
		if i == 0 && t.pinned != nil {
			annotation += fmt.Sprintf("[gray:-:-]%s Pinned Tweet[-:-:-]", icon.Pinned)
		}

		fmt.Fprintln(t.view, createTweetLayout(annotation, content, i, width))

		// 引用元ツイートを表示
		if quotedTweet != nil {
			if !appearance.HideQuoteTweetSeparator {
				fmt.Fprintln(t.view, createSeparator(appearance.QuoteTweetSeparator, width))
			}

			fmt.Fprintln(t.view, createTweetLayout("", quotedTweet, -1, width))
		}

		// セパレータを挿入しない
		if appearance.HideTweetSeparator {
			continue
		}

		// 末尾のツイート以外ならセパレータを挿入
		if i < t.GetTweetsCount()-1 {
			fmt.Fprintln(t.view, createSeparator(appearance.TweetSeparator, width))
		}
	}

	t.scrollToTweet(cursorPos)
}

// DrawMessage : ツイートビューにメッセージを表示
func (t *tweets) DrawMessage(s string) {
	t.view.Clear().
		SetTextAlign(tview.AlignCenter).
		SetText(s)
}

// scrollToTweet : 指定ツイートまでスクロール
func (t *tweets) scrollToTweet(i int) {
	// 範囲内に丸める
	if max := t.GetTweetsCount(); i < 0 {
		i = 0
	} else if i >= max {
		i = max - 1
	}

	t.view.Highlight(createTweetTag(i))
}

// moveCursor : カーソルを移動
func (t *tweets) moveCursor(c int) {
	idx := getHighlightId(t.view.GetHighlights())
	if idx == -1 {
		return
	}

	t.scrollToTweet(idx + int(c))
}
