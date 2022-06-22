package ui

import (
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
)

// timelineType : タイムラインの種類
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
		basePage: newBasePage(string(tt)),
		tlType:   tt,
	}

	page.SetFrame(page.tweets.view)
	page.frame.SetInputCapture(page.handleKeyEvents)

	return page
}

// Load : タイムライン読み込み
func (t *timelinePage) Load() {
	t.mu.Lock()
	defer t.mu.Unlock()

	var (
		tweets    []*twitter.TweetDictionary
		rateLimit *twitter.RateLimit
		err       error
	)

	shared.SetStatus(t.name, "Loading...")

	sinceID := t.tweets.GetSinceID()

	switch t.tlType {
	case homeTL:
		tweets, rateLimit, err = shared.api.FetchHomeTileline(shared.api.CurrentUser.ID, sinceID, 25)
	case mentionTL:
		tweets, rateLimit, err = shared.api.FetchUserMentionTimeline(shared.api.CurrentUser.ID, sinceID, 25)
	}

	if err != nil {
		t.tweets.DrawMessage(err.Error())
		shared.SetErrorStatus(t.name, err.Error())
		return
	}

	t.tweets.Register(tweets)
	t.tweets.Draw()

	t.updateIndicator("", rateLimit)
	t.updateLoadedStatus(len(tweets))
}

// handleKeyEvents : タイムラインページのキーハンドラ
func (t *timelinePage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	return handleCommonPageKeyEvent(t, event)
}
