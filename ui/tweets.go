package ui

import (
	"fmt"
	"sync"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tweets struct {
	textView *tview.TextView
	contents []*twitter.TweetDictionary
	count    int
	mu       sync.Mutex
}

func newTweets() *tweets {
	t := &tweets{
		textView: tview.NewTextView(),
		contents: []*twitter.TweetDictionary{},
		count:    0,
	}

	t.textView.SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true).
		SetRegions(true).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			t.textView.ScrollToHighlight()
		}).
		SetInputCapture(t.handleKeyEvents).
		SetBackgroundColor(tcell.ColorDefault)

	return t
}

func (t *tweets) getSinceID() string {
	if len(t.contents) == 0 {
		return ""
	}

	return t.contents[0].Tweet.ID
}

func (t *tweets) register(tweets []*twitter.TweetDictionary) {
	t.mu.Lock()
	defer t.mu.Unlock()

	length := len(tweets)
	t.contents = append(tweets, t.contents...)
	t.count = len(t.contents)

	shared.setStatus(fmt.Sprintf("%d tweets loaded", length))
}

func (t *tweets) draw() {
	width := getWindowWidth()

	t.textView.Clear()

	for i, content := range t.contents {
		var quotedTweet *twitter.TweetDictionary = nil

		// 参照ツイートを確認
		for _, rc := range content.ReferencedTweets {
			switch rc.Reference.Type {
			case "retweeted":
				fmt.Fprintln(t.textView, createAnnotation("RT by", content.Author))
				content = content.ReferencedTweets[0].TweetDictionary
			case "replied_to":
				fmt.Fprintln(t.textView, createAnnotation("Reply to", rc.TweetDictionary.Author))
			case "quoted":
				quotedTweet = rc.TweetDictionary
			}
		}

		// 表示部分を作成
		layout := createTweetLayout(content, i)
		fmt.Fprintln(t.textView, layout)

		// 引用元ツイートを表示
		if quotedTweet != nil {
			fmt.Fprintln(t.textView, createSeparator("-", width))
			layout := createTweetLayout(quotedTweet, -1)
			fmt.Fprintln(t.textView, layout)
		}

		// 末尾のツイートでないならセパレータを挿入
		if i < t.count-1 {
			fmt.Fprintln(t.textView, createSeparator("─", width))
		}
	}

	t.scrollToTweet(0)
}

func (t *tweets) scrollToTweet(i int) {
	t.textView.Highlight(createTweetId(i))
}

func (t *tweets) cursorUp() {
	idx := getHighlightId(t.textView.GetHighlights())
	if idx == -1 {
		return
	}

	if idx--; idx < 0 {
		idx = t.count - 1
	}

	t.scrollToTweet(idx)
}

func (t *tweets) cursorDown() {
	idx := getHighlightId(t.textView.GetHighlights())
	if idx == -1 {
		return
	}

	idx = (idx + 1) % t.count

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
		t.scrollToTweet(t.count - 1)
		return nil
	}

	return event
}
