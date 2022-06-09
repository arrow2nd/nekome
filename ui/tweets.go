package ui

import (
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/rivo/tview"
)

type tweets struct {
	textView *tview.TextView
	content  []*twitter.TweetObj
	index    int
}

func newTweets() *tweets {
	t := &tweets{
		textView: tview.NewTextView(),
		content:  []*twitter.TweetObj{},
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

func (t *tweets) draw() {
	fmt.Fprintln(t.textView, "test!")
}
