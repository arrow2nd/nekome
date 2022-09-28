package app

import (
	"fmt"
	"sync"

	"github.com/arrow2nd/nekome/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rivo/tview"
)

// cursorMove : カーソルの移動量
type cursorMove int

const (
	cursorMoveUp   cursorMove = -1
	cursorMoveDown cursorMove = 1
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

	t.view.SetHighlightedFunc(func(added, removed, remaining []string) {
		t.view.ScrollToHighlight()
	})

	if err := t.setKeyEventHandler(); err != nil {
		return nil, err
	}

	return t, nil
}

// setKeyEventHandler : キーハンドラを設定
func (t *tweets) setKeyEventHandler() error {
	handler := map[string]func(){
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
			t.actionForTweet(tweetLike)
		},
		config.ActionTweetUnlike: func() {
			t.actionForTweet(tweetUnlike)
		},
		config.ActionTweetRetweet: func() {
			t.actionForTweet(tweetRetweet)
		},
		config.ActionTweetUnretweet: func() {
			t.actionForTweet(tweetUnretweet)
		},
		config.ActionTweetRemove: func() {
			t.actionForTweet(tweetDelete)
		},
		config.ActionUserFollow: func() {
			t.actionForUser(userFollow)
		},
		config.ActionUserUnfollow: func() {
			t.actionForUser(userUnfollow)
		},
		config.ActionUserBlock: func() {
			t.actionForUser(userBlock)
		},
		config.ActionUserUnblock: func() {
			t.actionForUser(userUnblock)
		},
		config.ActionUserMute: func() {
			t.actionForUser(userMute)
		},
		config.ActionUserUnmute: func() {
			t.actionForUser(userUnmute)
		},
		config.ActionOpenUserPage: func() {
			t.openUserPage()
		},
		config.ActionQuote: func() {
			t.postQuoteTweet()
		},
		config.ActionReply: func() {
			t.postReply()
		},
		config.ActionOpenBrowser: func() {
			t.openBrower()
		},
		config.ActionCopyUrl: func() {
			t.copyLinkToClipBoard()
		},
	}

	c, err := shared.conf.Pref.Keybindings.Tweet.MappingEventHandler(handler)
	if err != nil {
		return err
	}

	t.view.SetInputCapture(c.Capture)

	return nil
}

// GetSinceID : 一番新しいツイートIDを取得
func (t *tweets) GetSinceID() string {
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
		// ピン留めツイート意外が選択されている
		c = t.contents[id-1]
	}

	// リツイートなら参照先に置き換え
	for _, rc := range c.ReferencedTweets {
		if rc.Reference.Type == "retweeted" {
			c = c.ReferencedTweets[0].TweetDictionary
		}
	}

	return c
}

// register : ツイートを登録
func (t *tweets) register(tweets []*twitter.TweetDictionary) int {
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

// RegisterPinned : ピン留めツイートを登録
func (t *tweets) RegisterPinned(tweet *twitter.TweetDictionary) {
	t.pinned = tweet
}

// UpdateRateLimit : レート制限を更新
func (t *tweets) UpdateRateLimit(r *twitter.RateLimit) {
	if r != nil {
		t.rateLimit = r
	}
}

// Update : ツイートを更新
func (t *tweets) Update(tweets []*twitter.TweetDictionary) {
	addedTweetsCount := t.register(tweets)

	// カーソル位置を決定
	cursorPos := getHighlightId(t.view.GetHighlights())
	if cursorPos == -1 {
		cursorPos = 0
	}

	// 先頭以外のツイートを選択中の場合、更新後もそのツイートを選択したままにする
	// NOTE: "先頭以外" なのは、ストリームモードで放置した時にカーソルが段々下に下がってしまうのを防ぐため
	if cursorPos != 0 {
		cursorPos += addedTweetsCount
	}

	t.draw(cursorPos)
}

// draw : 描画（表示幅はターミナルのウィンドウ幅に依存）
func (t *tweets) draw(cursorPos int) {
	width := getWindowWidth()

	// ビューの初期化
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

		// 参照ツイートを確認
		for _, rc := range content.ReferencedTweets {
			switch rc.Reference.Type {
			case "retweeted":
				fmt.Fprintln(t.view, createAnnotation("RT by", content.Author))
				content = content.ReferencedTweets[0].TweetDictionary
			case "replied_to":
				fmt.Fprintln(t.view, createAnnotation("Reply to", rc.TweetDictionary.Author))
			case "quoted":
				quotedTweet = rc.TweetDictionary
			}
		}

		// ピン留めツイート
		if i == 0 && t.pinned != nil {
			icon := shared.conf.Pref.Icon.Pinned
			fmt.Fprintf(t.view, "[gray:-:-]%s Pinned Tweet[-:-:-]\n", icon)
		}

		// 表示部分を作成
		layout := createTweetLayout(content, i, width)
		fmt.Fprintln(t.view, layout)

		// 引用元ツイートを表示
		if quotedTweet != nil {
			fmt.Fprintln(t.view, createSeparator("-", width))
			fmt.Fprintln(t.view, createTweetLayout(quotedTweet, -1, width))
		}

		// 末尾のツイートでないならセパレータを挿入
		lastIndex := t.GetTweetsCount() - 1
		if i < lastIndex {
			fmt.Fprintln(t.view, createSeparator("─", width))
		}
	}

	t.scrollToTweet(cursorPos)
}

// DrawMessage : Viewにメッセージを表示
func (t *tweets) DrawMessage(s string) {
	t.view.Clear().
		SetTextAlign(tview.AlignCenter).
		SetText(s)
}

// scrollToTweet : 指定ツイートまでスクロール
func (t *tweets) scrollToTweet(i int) {
	// 範囲内に丸める
	if max := t.GetTweetsCount(); i < 0 {
		i = max - 1
	} else if i >= max {
		i = 0
	}

	t.view.Highlight(createTweetTag(i))
}

// moveCursor : カーソルを移動
func (t *tweets) moveCursor(c cursorMove) {
	idx := getHighlightId(t.view.GetHighlights())
	if idx == -1 {
		return
	}

	t.scrollToTweet(idx + int(c))
}
