package ui

import (
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// timelineType タイムラインの種類
type timelineType string

const (
	homeTL    timelineType = "Home"
	mentionTL timelineType = "Mention"
)

type timelinePage struct {
	*basePage
	tlType timelineType
}

func newTimelinePage(tt timelineType) *timelinePage {
	page := &timelinePage{
		basePage: newBasePage(),
		tlType:   tt,
	}

	page.frame = tview.NewFrame(page.tweets.textView).
		SetBorders(1, 1, 0, 0, 1, 1)

	page.frame.SetInputCapture(page.handleKeyEvents)

	return page
}

// Load タイムライン読み込み
func (t *timelinePage) Load() {
	var (
		tweets []*twitter.TweetDictionary
		err    error
	)

	shared.setStatus("Loading...")

	sinceID := t.tweets.getSinceID()

	switch t.tlType {
	case homeTL:
		tweets, err = shared.api.FetchHomeTileline(shared.api.CurrentUser.ID, sinceID, 50)
	case mentionTL:
		tweets, err = shared.api.FetchUserMentionTimeline(shared.api.CurrentUser.ID, sinceID, 25)
	}

	if err != nil {
		shared.setStatus(err.Error())
		return
	}

	count := t.tweets.register(tweets)
	t.tweets.draw()

	shared.setStatus(fmt.Sprintf("[%s] %d tweets loaded", t.tlType, count))
}

func (t *timelinePage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handlePageKeyEvents(t, event)
}
