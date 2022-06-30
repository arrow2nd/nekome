package app

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
	*tweetsBasePage
	tlType timelineType
}

// newTimelinePage : タイムラインページを作成
func newTimelinePage(tt timelineType) *timelinePage {
	tabName := shared.conf.Settings.Texts.TabHome
	if tt == mentionTL {
		tabName = shared.conf.Settings.Texts.TabMention
	}

	page := &timelinePage{
		tweetsBasePage: newTweetsBasePage(tabName),
		tlType:         tt,
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

	shared.SetStatus(t.name, shared.conf.Settings.Texts.Loading)

	// タイムラインを取得
	id := shared.api.CurrentUser.ID
	count := shared.conf.Settings.Feature.LoadTweetsCount
	sinceID := t.tweets.GetSinceID()

	if t.tlType == homeTL {
		tweets, rateLimit, err = shared.api.FetchHomeTileline(id, sinceID, count)
	} else {
		tweets, rateLimit, err = shared.api.FetchUserMentionTimeline(id, sinceID, count)
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
