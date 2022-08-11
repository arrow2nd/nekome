package app

import (
	"context"
	"time"

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
	cancel context.CancelFunc
}

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
func (t *timelinePage) Load(focus bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var (
		tweets    []*twitter.TweetDictionary
		rateLimit *twitter.RateLimit
		err       error
	)

	// 読み込み中表示
	if !t.isStreamMode() {
		shared.SetStatus(t.name, shared.conf.Settings.Texts.Loading)
	}

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

	t.tweets.Register(tweets, rateLimit)
	t.tweets.Draw()

	t.updateIndicator(t.getStreamStatus(), focus)

	// 読み込み完了表示
	if !t.isStreamMode() {
		t.updateLoadedStatus(len(tweets))
	}
}

// getStreamStatus : ストリームモードのステータスを取得
func (t *timelinePage) getStreamStatus() string {
	if t.isStreamMode() {
		return "Stream Mode | "
	}

	return ""
}

// isStreamMode : ストリームモードが有効かどうか
func (t *timelinePage) isStreamMode() bool {
	return t.cancel != nil
}

// startStream : ストリームモード開始
func (t *timelinePage) startStream() {
	if t.isStreamMode() {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go t.stream(ctx, (15*60)/150)

	t.updateIndicator(t.getStreamStatus(), true)
}

// stopStream : ストリームモード停止
func (t *timelinePage) stopStream() {
	if !t.isStreamMode() {
		return
	}

	t.cancel()
	t.cancel = nil

	t.updateIndicator(t.getStreamStatus(), true)
}

// stream : ストリームモード
func (t *timelinePage) stream(ctx context.Context, intervalSec time.Duration) {
	ticker := time.NewTicker(intervalSec * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.Load(true)
		}
	}
}

// handleKeyEvents : タイムラインページのキーハンドラ
func (t *timelinePage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	if keyRune == 's' {
		t.startStream()
		return nil
	}

	if keyRune == 'S' {
		t.stopStream()
		return nil
	}

	return handleCommonPageKeyEvent(t, event)
}
