package ui

import (
	"fmt"
	"sync"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tweets struct {
	view          *tview.TextView
	pinned        *twitter.TweetDictionary
	contents      []*twitter.TweetDictionary
	contentsCount int
	mu            sync.Mutex
}

func newTweets() *tweets {
	t := &tweets{
		view:          tview.NewTextView(),
		pinned:        nil,
		contents:      []*twitter.TweetDictionary{},
		contentsCount: 0,
	}

	t.view.SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			t.view.ScrollToHighlight()
		}).
		SetInputCapture(t.handleKeyEvents).
		SetBackgroundColor(tcell.ColorDefault)

	return t
}

// GetSinceID : 一番新しいツイートIDを取得
func (t *tweets) GetSinceID() string {
	if len(t.contents) == 0 {
		return ""
	}

	return t.contents[0].Tweet.ID
}

// Register : ツイートを登録
func (t *tweets) Register(tweets []*twitter.TweetDictionary) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.contents = append(tweets, t.contents...)
}

// RegisterPinned : ピン留めツイートを登録
func (t *tweets) RegisterPinned(tweet *twitter.TweetDictionary) {
	t.pinned = tweet
}

// Draw : 描画（表示幅はターミナルのウィンドウ幅に依存）
func (t *tweets) Draw() {
	width := getWindowWidth()

	t.view.Clear()

	contents := t.contents

	// ピン留めツイートがある場合、先頭に追加
	if t.pinned != nil {
		contents = append([]*twitter.TweetDictionary{t.pinned}, t.contents...)
	}

	t.contentsCount = len(contents)

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
			fmt.Fprintln(t.view, "[gray:-:-] Pinned Tweet[r]")
		}

		// 表示部分を作成
		layout := createTweetLayout(content, i)
		fmt.Fprintln(t.view, layout)

		// 引用元ツイートを表示
		if quotedTweet != nil {
			fmt.Fprintln(t.view, createSeparator("-", width))
			layout := createTweetLayout(quotedTweet, -1)
			fmt.Fprintln(t.view, layout)
		}

		// 末尾のツイートでないならセパレータを挿入
		if i < t.contentsCount-1 {
			fmt.Fprintln(t.view, createSeparator("─", width))
		}
	}

	t.scrollToTweet(0)
}

func (t *tweets) scrollToTweet(i int) {
	t.view.Highlight(createTweetId(i))
}

func (t *tweets) cursorUp() {
	idx := getHighlightId(t.view.GetHighlights())
	if idx == -1 {
		return
	}

	if idx--; idx < 0 {
		idx = t.contentsCount - 1
	}

	t.scrollToTweet(idx)
}

func (t *tweets) cursorDown() {
	idx := getHighlightId(t.view.GetHighlights())
	if idx == -1 {
		return
	}

	idx = (idx + 1) % t.contentsCount

	t.scrollToTweet(idx)
}

func (t *tweets) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	keyRune := event.Rune()

	if key == tcell.KeyUp || keyRune == 'k' {
		t.cursorUp()
		return nil
	}

	if key == tcell.KeyDown || keyRune == 'j' {
		t.cursorDown()
		return nil
	}

	if keyRune == 'g' {
		t.scrollToTweet(0)
		return nil
	}

	if keyRune == 'G' {
		t.scrollToTweet(t.contentsCount - 1)
		return nil
	}

	return event
}
