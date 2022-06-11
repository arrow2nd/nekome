package ui

import (
	"fmt"
	"sync"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rivo/tview"
)

type tweets struct {
	textView *tview.TextView
	contents []*twitter.TweetObj
	index    int
	mu       sync.Mutex
}

func newTweets() *tweets {
	t := &tweets{
		textView: tview.NewTextView(),
		contents: []*twitter.TweetObj{},
		index:    0,
	}

	t.textView.SetDynamicColors(true).
		SetScrollable(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			t.textView.ScrollToHighlight()
		})

	return t
}

func (t *tweets) getSinceID() string {
	if len(t.contents) == 0 {
		return ""
	}

	return t.contents[0].ID
}

func (t *tweets) register(tweets []*twitter.TweetObj) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.index = 0
	length := len(tweets)
	t.contents = append(tweets, t.contents...)

	shared.setStatus(fmt.Sprintf("%d tweets loaded", length))
}

func (t *tweets) draw() {
	t.textView.Clear()

	for _, tweet := range t.contents {
		fmt.Fprintf(t.textView, "%s\n\n", tweet.Text)
	}
}
