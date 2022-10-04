package app

import (
	"fmt"
	"sync"

	"github.com/arrow2nd/nekome/v2/config"
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

// createCommonPageKeyHandler : ページ共通のキーハンドラを作成
func createCommonPageKeyHandler(p page) (func(*tcell.EventKey) *tcell.EventKey, error) {
	handler := map[string]func(){
		config.ActionReloadPage: func() {
			go p.Load()
		},
	}

	c, err := shared.conf.Pref.Keybindings.Page.MappingEventHandler(handler)
	if err != nil {
		return nil, err
	}

	return c.Capture, nil
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
		name:      truncate(name, shared.conf.Pref.Appearance.TabMaxWidth),
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

func newTweetsBasePage(name string) (*tweetsBasePage, error) {
	tweets, err := newTweets()
	if err != nil {
		return nil, err
	}

	return &tweetsBasePage{
		basePage: newBasePage(name),
		tweets:   tweets,
	}, nil
}

// updateIndicator : インジケータを更新
func (t *tweetsBasePage) updateIndicator(s string) {
	rateLimit := t.tweets.rateLimit

	// APIリミット
	apiLimit := "unknown"
	if rateLimit != nil {
		apiLimit = fmt.Sprintf("%d / %d", rateLimit.Remaining, rateLimit.Limit)
	}

	t.indicator = s + fmt.Sprintf("API limit: %s", apiLimit)

	// ページがアクティブなら表示を更新する
	if t.isActive {
		shared.SetIndicator(t.indicator)
	}
}

// updateLoadedStatus : ステータスメッセージを更新
func (t *tweetsBasePage) updateLoadedStatus(count int) {
	text := shared.conf.Pref.Text.NoTweets

	if count > 0 {
		text = fmt.Sprintf("%d tweets loaded", count)
	}

	shared.SetStatus(t.name, text)
}
