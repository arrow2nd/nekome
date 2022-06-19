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

func handlePageKeyEvents(p page, event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	// リロード
	if keyRune == 'R' {
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

// GetName ページ名を取得
func (b *basePage) GetName() string {
	return b.name
}

// SetFrame フレームを設定
func (b *basePage) SetFrame(p tview.Primitive) {
	b.frame = tview.NewFrame(p)
	b.frame.SetBorders(1, 1, 0, 0, 1, 1)
}

// GetPrimivite プリミティブを取得
func (b *basePage) GetPrimivite() tview.Primitive {
	return b.frame
}

func (b *basePage) showLoadedStatus(r *twitter.RateLimit) {
	text := fmt.Sprintf("%d tweets loaded (API limit: %d / %d)", b.tweets.count, r.Remaining, r.Limit)
	shared.setStatus(b.name, text)
}
