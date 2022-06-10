package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type homeTimeline struct {
	frame  *tview.Frame
	tweets *tweets
}

func newHomeTimeline() *homeTimeline {
	home := &homeTimeline{
		frame:  nil,
		tweets: newTweets(),
	}

	home.frame = tview.NewFrame(home.tweets.textView).
		SetBorders(0, 0, 0, 0, 1, 1)

	home.setHomeTimelineKeyEvents()

	return home
}

func (h *homeTimeline) init() {
	h.tweets.draw()
}

func (h *homeTimeline) setHomeTimelineKeyEvents() {
	h.frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'R':
			setStatus("Reload!")
			return nil
		case 'Q':
			setStatus("Quit!")
			return nil
		}

		return event
	})
}
