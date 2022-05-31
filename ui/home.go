package ui

import "github.com/rivo/tview"

type homeTimeline struct {
	frame  *tview.Frame
	tweets *tweets
}

func newHomeTimeline() *homeTimeline {
	home := homeTimeline{
		frame:  nil,
		tweets: newTweets(),
	}

	home.frame = tview.NewFrame(home.tweets.textView).
		SetBorders(0, 0, 0, 0, 1, 1)

	return &home
}

func (h *homeTimeline) init() {
	h.tweets.draw()
}
