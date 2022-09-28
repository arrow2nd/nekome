package app

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/arrow2nd/nekome/config"
	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
)

// 読み込み間隔
const (
	reloadIntervalMin     = 5
	reloadIntervalDefault = 10
)

// timelineType : タイムラインの種類
type timelineType string

const (
	homeTimeline    timelineType = "Home"
	mentionTimeline timelineType = "Mention"
)

type timelinePage struct {
	*tweetsBasePage
	tlType         timelineType
	reloadInterval time.Duration
	cancel         context.CancelFunc
}

func newTimelinePage(t timelineType) (*timelinePage, error) {
	tabName := shared.conf.Pref.Text.TabHome
	if t == mentionTimeline {
		tabName = shared.conf.Pref.Text.TabMention
	}

	basePage, err := newTweetsBasePage(tabName)
	if err != nil {
		return nil, err
	}

	p := &timelinePage{
		tweetsBasePage: basePage,
		tlType:         t,
		reloadInterval: 0,
		cancel:         nil,
	}

	p.SetFrame(p.tweets.view)

	if err := p.setKeyHandler(); err != nil {
		return nil, err
	}

	return p, nil
}

// setKeyHandler : キーハンドラを設定
func (t *timelinePage) setKeyHandler() error {
	handler := map[string]func(){
		config.ActionStreamModeStart: func() {
			t.startStream()
		},
		config.ActionStreamModeStop: func() {
			t.closeStream()
		},
	}

	c, err := shared.conf.Pref.Keybindings.HomeTimeline.MappingEventHandler(handler)
	if err != nil {
		return err
	}

	commonHandler, err := createCommonPageKeyHandler(t)
	if err != nil {
		return err
	}

	t.frame.SetInputCapture(func(ev *tcell.EventKey) *tcell.EventKey {
		if e := c.Capture(ev); e == nil {
			return nil
		}

		// ストリームモード中は共通のキーハンドラを無効化（手動リロード禁止）
		if t.isStreamMode() {
			return nil
		}

		return commonHandler(ev)
	})

	return nil
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
		shared.SetStatus(t.name, shared.conf.Pref.Text.Loading)
	}

	// タイムラインを取得
	id := shared.api.CurrentUser.ID
	count := shared.conf.Pref.Feature.LoadTweetsLimit
	sinceID := t.tweets.GetSinceID()

	if t.tlType == homeTimeline {
		tweets, rateLimit, err = shared.api.FetchHomeTileline(id, sinceID, count)
	} else {
		tweets, rateLimit, err = shared.api.FetchUserMentionTimeline(id, sinceID, count)
	}

	if err != nil {
		t.tweets.DrawMessage(err.Error())
		shared.SetErrorStatus(t.name, "failed to retrieve timeline")
		return
	}

	t.tweets.Update(tweets)
	t.tweets.UpdateRateLimit(rateLimit)

	t.updateIndicator(t.getStreamIndicator())

	// 読み込み完了表示
	if !t.isStreamMode() {
		t.updateLoadedStatus(len(tweets))
	}
}

// OnDelete : ページが破棄された
func (t *timelinePage) OnDelete() {
	// ストリームモードが有効なら終了する
	if t.isStreamMode() {
		t.closeStream()
	}
}

// getStreamIndicator : ストリームモードのインジケータを取得
func (t *timelinePage) getStreamIndicator() string {
	if !t.isStreamMode() {
		return ""
	}

	return fmt.Sprintf("Stream Mode | Interval: %ds | ", t.reloadInterval)
}

// isStreamMode : ストリームモードが有効かどうか
func (t *timelinePage) isStreamMode() bool {
	return t.cancel != nil
}

// startStream : ストリームモードを開始
func (t *timelinePage) startStream() {
	if t.isStreamMode() {
		shared.SetErrorStatus(t.name, "stream mode has already started")
		return
	}

	// 読み込み間隔を決定
	if interval, err := t.calcReloadInterval(); err != nil {
		shared.SetErrorStatus(t.name, fmt.Sprintf("stream mode cannot be started (%s)", err.Error()))
		return
	} else {
		t.reloadInterval = interval
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel

	go t.stream(ctx)

	shared.SetStatus(t.name, "stream mode has been started")
	t.updateIndicator(t.getStreamIndicator())
}

// closeStream : ストリームモードを終了
func (t *timelinePage) closeStream() {
	if !t.isStreamMode() {
		shared.SetErrorStatus(t.name, "stream mode has not been started")
		return
	}

	t.cancel()
	t.cancel = nil

	shared.SetStatus(t.name, "stream mode has been closed")
	t.updateIndicator(t.getStreamIndicator())
}

// calcReloadInterval : 読み込み間隔を計算
func (t *timelinePage) calcReloadInterval() (time.Duration, error) {
	// レート制限が取得できない
	if t.tweets.rateLimit == nil {
		return 0, errors.New("failed to obtain rate limit")
	}

	remainingSec := time.Until(t.tweets.rateLimit.Reset.Time()).Seconds()

	// レート制限を超えている場合、デフォルト値を返す
	if remainingSec <= 0 || t.tweets.rateLimit.Remaining <= 0 {
		return reloadIntervalDefault, nil
	}

	newInterval := math.Round(remainingSec / float64(t.tweets.rateLimit.Remaining))

	// 最小間隔は5秒
	if newInterval < reloadIntervalMin {
		return reloadIntervalMin, nil
	}

	return time.Duration(newInterval), nil
}

// stream : ストリームモード
func (t *timelinePage) stream(ctx context.Context) {
	ticker := time.NewTicker(t.reloadInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.loadStream(ticker)
		}
	}
}

// loadStream : ストリームを読み込み
func (t *timelinePage) loadStream(ticker *time.Ticker) {
	// レート制限情報が無い場合は不正な状態なので終了させる
	if t.tweets.rateLimit == nil {
		t.closeStream()
		shared.SetErrorStatus(t.name, "stream mode has been interrupted (failed to obtain rate limit)")
		return
	}

	prevRemaining := t.tweets.rateLimit.Remaining

	t.Load()

	if t.tweets.rateLimit.Remaining <= prevRemaining {
		return
	}

	// レート制限がリセットされたら、読み込み間隔を再計算する
	if nextInterval, _ := t.calcReloadInterval(); t.reloadInterval != nextInterval {
		t.reloadInterval = nextInterval
		ticker.Reset(nextInterval * time.Second)
	}
}
