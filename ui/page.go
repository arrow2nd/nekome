package ui

import (
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
