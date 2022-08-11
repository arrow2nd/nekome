package app

import (
	"context"
	"fmt"
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
func (t *timelinePage) Load() {
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

	t.tweets.Register(tweets)
	t.tweets.UpdateRateLimit(rateLimit)
	t.tweets.Draw()

	t.updateIndicator(t.getStreamIndicator())

	// 読み込み完了表示
	if !t.isStreamMode() {
		t.updateLoadedStatus(len(tweets))
	}
}

// getStreamIndicator : ストリームモードのインジケータを取得
func (t *timelinePage) getStreamIndicator() string {
	if !t.isStreamMode() {
		return ""
	}

	return "Stream Mode | "
}

// isStreamMode : ストリームモードが有効かどうか
func (t *timelinePage) isStreamMode() bool {
	return t.cancel != nil
}

// startStream : ストリームモード開始
func (t *timelinePage) startStream() {
	if t.isStreamMode() {
		shared.SetErrorStatus(t.name, "stream mode has already started")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	// 更新間隔を決定
	nextRefreshSpan := t.tweets.rateLimit.Reset.Time().Unix() - time.Now().Unix()
	bufferSec := nextRefreshSpan / int64(t.tweets.rateLimit.Remaining)
	if bufferSec < 5 {
		bufferSec = 5
	}

	shared.SetStatus(t.name, fmt.Sprintf("Span: %ds", bufferSec))

	go t.stream(ctx, time.Duration(bufferSec))

	t.updateIndicator(t.getStreamIndicator())
}

// endStream : ストリームモード終了
func (t *timelinePage) endStream() {
	if !t.isStreamMode() {
		shared.SetErrorStatus(t.name, "stream mode has not been started")
		return
	}

	t.cancel()
	t.cancel = nil

	t.updateIndicator(t.getStreamIndicator())
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
			t.Load()
		}
	}
}

// handleKeyEvents : タイムラインページのキーハンドラ
func (t *timelinePage) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	// ストリームモード開始
	if keyRune == 's' {
		t.startStream()
		return nil
	}

	// ストリームモード終了
	if keyRune == 'S' {
		t.endStream()
		return nil
	}

	return handleCommonPageKeyEvent(t, event)
}
