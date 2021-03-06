package app

import (
	"fmt"
	"sync"

	"github.com/g8rswimmer/go-twitter/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type page interface {
	GetName() string
	GetPrimivite() tview.Primitive
	Load()
	OnVisible()
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
}

func newBasePage(name string) *basePage {
	return &basePage{
		name:      truncate(name, shared.conf.Settings.Appearance.TabMaxWidth),
		indicator: "",
		frame:     nil,
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

// OnVisible : ページが表示された際に呼ばれるコールバック
func (b *basePage) OnVisible() {
	// 以前のインジケータの内容を反映
	shared.SetIndicator(b.indicator)
}

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
func (t *tweetsBasePage) updateIndicator(s string, r *twitter.RateLimit) {
	t.indicator = fmt.Sprintf("%sAPI limit: %d / %d", s, r.Remaining, r.Limit)
	shared.SetIndicator(t.indicator)
}

// updateLoadedStatus : ステータスメッセージを更新
func (t *tweetsBasePage) updateLoadedStatus(count int) {
	text := shared.conf.Settings.Texts.NoTweets

	if count > 0 {
		text = fmt.Sprintf("%d tweets loaded", count)
	}

	shared.SetStatus(t.name, text)
}
