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

// タイムラインのタイプ
const (
	timelineTypeHome    string = "Home"
	timelineTypeMention string = "Mention"
)

type timelinePage struct {
	*tweetsBasePage
	timelineType   string
	reloadInterval time.Duration
	cancel         context.CancelFunc
}

func newTimelinePage(t string) (*timelinePage, error) {
	tabName := shared.conf.Pref.Text.TabHome
	if t == timelineTypeMention {
		tabName = shared.conf.Pref.Text.TabMention
	}

	basePage, err := newTweetsBasePage(tabName)
	if err != nil {
		return nil, err
	}

	p := &timelinePage{
		tweetsBasePage: basePage,
		timelineType:   t,
		reloadInterval: 0,
		cancel:         nil,
	}

	p.SetFrame(p.tweets.view)

	if err := p.setKeybindings(); err != nil {
		return nil, err
	}

	return p, nil
}

// setKeybindings : キーバインドを設定
func (t *timelinePage) setKeybindings() error {
	handlers := map[string]func(){
		config.ActionStreamModeStart: func() {
			t.startStreamMode()
		},
		config.ActionStreamModeStop: func() {
			t.stopStreamMode()
		},
	}

	c, err := shared.conf.Pref.Keybindings.HomeTimeline.MappingEventHandler(handlers)
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

		// ストリームモード中はページ共通のキーバインドを無効化（手動リロード禁止）
		if t.isStreamMode() {
			return nil
		}

		return commonHandler(ev)
	})

	return nil
}

// Load : タイムライン読み込み
func (t *timelinePage) Load() {
	var (
		tweets    []*twitter.TweetDictionary
		rateLimit *twitter.RateLimit
		err       error
	)

	t.mu.Lock()
	defer t.mu.Unlock()

	// 読み込み中表示
	if !t.isStreamMode() {
		shared.SetStatus(t.name, shared.conf.Pref.Text.Loading)
	}

	id := shared.api.CurrentUser.ID
	count := shared.conf.Pref.Feature.LoadTweetsLimit
	sinceId := t.tweets.GetSinceID()

	// タイムラインを取得
	switch t.timelineType {
	case timelineTypeHome:
		tweets, rateLimit, err = shared.api.FetchHomeTileline(id, sinceId, count)
	case timelineTypeMention:
		tweets, rateLimit, err = shared.api.FetchUserMentionTimeline(id, sinceId, count)
	default:
		t.tweets.DrawMessage("Load error")
		shared.SetErrorStatus(t.name, fmt.Sprintf("unknown timeline type: %s", t.timelineType))
		return
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
	// ストリームモード実行中なら停止する
	if t.isStreamMode() {
		t.stopStreamMode()
	}
}

// getStreamIndicator : ストリームモードのインジケータを取得
func (t *timelinePage) getStreamIndicator() string {
	if !t.isStreamMode() {
		return ""
	}

	return fmt.Sprintf("Stream Mode | Interval: %ds | ", t.reloadInterval)
}

// isStreamMode : ストリームモードが実行中かどうか
func (t *timelinePage) isStreamMode() bool {
	return t.cancel != nil
}

// startStreamMode : ストリームモードを開始
func (t *timelinePage) startStreamMode() {
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

	go t.streamMainLoop(ctx)

	shared.SetStatus(t.name, "stream mode has been started")
	t.updateIndicator(t.getStreamIndicator())
}

// stopStreamMode : ストリームモードを停止
func (t *timelinePage) stopStreamMode() {
	if !t.isStreamMode() {
		shared.SetErrorStatus(t.name, "stream mode has not been started")
		return
	}

	t.cancel()
	t.cancel = nil

	shared.SetStatus(t.name, "stream mode has been closed")
	t.updateIndicator(t.getStreamIndicator())
}

// streamMainLoop : ストリームモードメインループ
func (t *timelinePage) streamMainLoop(ctx context.Context) {
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
		t.stopStreamMode()
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
