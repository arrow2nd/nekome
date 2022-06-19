package ui

import (
	"fmt"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
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

	page.SetFrame(page.tweets.view)
	page.frame.SetInputCapture(page.handleKeyEvents)

	return page
}

// Load タイムライン読み込み
func (t *timelinePage) Load() {
	t.mu.Lock()
	defer t.mu.Unlock()

	var (
		tweets []*twitter.TweetDictionary
		err    error
		label  string = string(t.tlType)
	)

	shared.setStatus(label, "Loading...")

	sinceID := t.tweets.getSinceID()

	switch t.tlType {
	case homeTL:
		tweets, err = shared.api.FetchHomeTileline(shared.api.CurrentUser.ID, sinceID, 25)
	case mentionTL:
		tweets, err = shared.api.FetchUserMentionTimeline(shared.api.CurrentUser.ID, sinceID, 25)
	}

	if err != nil {
		shared.setErrorStatus(label, err.Error())
		return
	}

	count := t.tweets.register(tweets)
	t.tweets.draw()

	shared.setStatus(label, fmt.Sprintf("%d tweets loaded", count))
}

func (t *timelinePage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handlePageKeyEvents(t, event)
}
