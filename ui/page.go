package ui

import (
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type page interface {
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
	frame  *tview.Frame
	tweets *tweets
	mu     sync.Mutex
}

func newBasePage() *basePage {
	return &basePage{
		frame:  nil,
		tweets: newTweets(),
	}
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
