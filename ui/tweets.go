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
	textView := tview.NewTextView()

	textView.SetDynamicColors(true).
		SetScrollable(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			textView.ScrollToHighlight()
		})

	return &tweets{
		textView: textView,
		content:  []*twitter.TweetObj{},
		index:    0,
	}
}

func (t *tweets) draw() {
	fmt.Fprintln(t.textView, "test!")
}
