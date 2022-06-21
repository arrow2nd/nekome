package ui

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
}

// handlePageKeyEvents : ページの共通キーハンドラ
func handlePageKeyEvents(p page, event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	// リロード
	if keyRune == 'r' {
		go p.Load()
		return nil
	}

	return event
}

type basePage struct {
	page
	name   string
	frame  *tview.Frame
	tweets *tweets
	mu     sync.Mutex
}

func newBasePage(name string) *basePage {
	return &basePage{
		name:   name,
		frame:  nil,
		tweets: newTweets(),
	}
}

// GetName : ページ名を取得
func (b *basePage) GetName() string {
	return b.name
}

// SetFrame : フレームを設定
func (b *basePage) SetFrame(p tview.Primitive) {
	b.frame = tview.NewFrame(p)
	b.frame.SetBorders(1, 1, 0, 0, 1, 1)
}

// GetPrimivite : プリミティブを取得
func (b *basePage) GetPrimivite() tview.Primitive {
	return b.frame
}

// showLoadedStatus : ロード後のステータスメッセージを設定
func (b *basePage) showLoadedStatus(count int, r *twitter.RateLimit) {
	text := ""
	limit := fmt.Sprintf("(API limit: %d / %d)", r.Remaining, r.Limit)

	if count <= 0 {
		text = "no new tweets " + limit
	} else {
		text = fmt.Sprintf("%d tweets loaded %s", count, limit)
	}

	shared.SetStatus(b.name, text)
}
