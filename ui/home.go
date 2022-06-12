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

func (h *homeTimeline) load() {
	shared.setStatus("Loading...")

	sinceID := h.tweets.getSinceID()
	tweets, err := shared.api.FetchHomeTileline(shared.api.CurrentUser.ID, sinceID, 100)
	if err != nil {
		shared.setStatus(err.Error())
		return
	}

	h.tweets.register(tweets)
	h.tweets.draw()
}

func (h *homeTimeline) setHomeTimelineKeyEvents() {
	h.frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'R':
				h.load()
				return nil
			}
		}

		return event
	})
}
