package app

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type page interface {
	GetName() string
	GetPrimivite() tview.Primitive
	Load()
	OnActive()
	OnInactive()
	OnDelete()
}

// handleCommonPageKeyEvent : ページ共通のキーハンドラ
func handleCommonPageKeyEvent(p page, event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	// リロード
	if keyRune == '.' {
		go p.Load()
		return nil
	}

	return event
}

type basePage struct {
	page
	name      string
	indicator string
	frame     *tview.Frame
	isActive  bool
}

func newBasePage(name string) *basePage {
	return &basePage{
		name:      truncate(name, shared.conf.Settings.Appearance.TabMaxWidth),
		indicator: "",
		frame:     nil,
		isActive:  false,
	}
}

// GetName : ページ名を取得
func (b *basePage) GetName() string {
	return b.name
}

// GetPrimivite : プリミティブを取得
func (b *basePage) GetPrimivite() tview.Primitive {
	return b.frame
}

// SetFrame : フレームを設定
func (b *basePage) SetFrame(p tview.Primitive) {
	b.frame = tview.NewFrame(p)
	b.frame.SetBorders(1, 1, 0, 0, 1, 1)
}

// Load : 読み込み
func (b *basePage) Load() {}

// OnActive : ページがアクティブになった
func (b *basePage) OnActive() {
	b.isActive = true

	// 以前のインジケータの内容を反映
	shared.SetIndicator(b.indicator)
}

// OnInactive : ページが非アクティブになった
func (b *basePage) OnInactive() {
	b.isActive = false
}

// OnDelete : ページが破棄された
func (b *basePage) OnDelete() {}

type tweetsBasePage struct {
	*basePage
	tweets *tweets
	mu     sync.Mutex
}

func newTweetsBasePage(name string) *tweetsBasePage {
	return &tweetsBasePage{
		basePage: newBasePage(name),
		tweets:   newTweets(),
	}
}

// updateIndicator : インジケータを更新
func (t *tweetsBasePage) updateIndicator(s string) {
	// APIリミット
	apiLimit := "unknown"
	if t.tweets.rateLimit != nil {
		apiLimit = fmt.Sprintf("%d / %d", t.tweets.rateLimit.Remaining, t.tweets.rateLimit.Limit)
	}

	t.indicator = s + fmt.Sprintf("API limit: %s", apiLimit)

	// ページがアクティブなら表示を更新する
	if t.isActive {
		shared.SetIndicator(t.indicator)
	}
}

// updateLoadedStatus : ステータスメッセージを更新
func (t *tweetsBasePage) updateLoadedStatus(count int) {
	text := shared.conf.Settings.Texts.NoTweets

	if count > 0 {
		text = fmt.Sprintf("%d tweets loaded", count)
	}

	shared.SetStatus(t.name, text)
}
