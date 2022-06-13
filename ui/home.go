package ui

import (
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type timelineType string

const (
	HomeTL    timelineType = "Home"
	MentionTL timelineType = "Mention"
)

type timelinePage struct {
	tlType timelineType
	frame  *tview.Frame
	tweets *tweets
}

func newTimelinePage(tl timelineType) *timelinePage {
	home := &timelinePage{
		tlType: tl,
		frame:  nil,
		tweets: newTweets(),
	}

	home.frame = tview.NewFrame(home.tweets.textView).
		SetBorders(1, 1, 0, 0, 1, 1)

	home.frame.SetInputCapture(home.handleTimelinePageKeyEvents)

	return home
}

func (t *timelinePage) GetPrimivite() tview.Primitive {
	return t.frame
}

func (t *timelinePage) Load() {
	var (
		tweets []*twitter.TweetDictionary
		err    error
	)

	defer shared.drawApplication()

	shared.setStatus("Loading...")

	sinceID := t.tweets.getSinceID()

	switch t.tlType {
	case HomeTL:
		tweets, err = shared.api.FetchHomeTileline(shared.api.CurrentUser.ID, sinceID, 25)
	case MentionTL:
		tweets, err = shared.api.FetchUserMentionTimeline(shared.api.CurrentUser.ID, sinceID, 25)
	}

	if err != nil {
		shared.setStatus(err.Error())
		return
	}

	t.tweets.register(tweets)
	t.tweets.draw()
}

func (t *timelinePage) handleTimelinePageKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'R':
			go t.Load()
			return nil
		}
	}

	return event
}
