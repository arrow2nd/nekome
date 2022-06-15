package ui

import (
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
	tlType timelineType
	frame  *tview.Frame
	tweets *tweets
}

func newTimelinePage(tl timelineType) *timelinePage {
	page := &timelinePage{
		tlType: tl,
		frame:  nil,
		tweets: newTweets(),
	}

	page.frame = tview.NewFrame(page.tweets.textView).
		SetBorders(1, 1, 0, 0, 1, 1)

	page.frame.SetInputCapture(page.handleTimelinePageKeyEvents)

	return page
}

// GetPrimivite プリミティブを取得
func (t *timelinePage) GetPrimivite() tview.Primitive {
	return t.frame
}

// Load タイムライン読み込み
func (t *timelinePage) Load() {
	var (
		tweets []*twitter.TweetDictionary
		err    error
	)

	defer shared.reqestDrawApp()

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

	t.tweets.register(tweets)
	t.tweets.draw()
}

func (t *timelinePage) handleTimelinePageKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	// リロード
	if keyRune == 'R' {
		go t.Load()
		return nil
	}

	return event
}
