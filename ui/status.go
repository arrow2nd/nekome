package ui

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type status struct {
	textView *tview.TextView
	text     string
	mu       sync.Mutex
}

func newStatus() *status {
	s := &status{
		textView: tview.NewTextView(),
		text:     "",
	}

	s.textView.SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDarkCyan)

	return s
}

func (s *status) set(t string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.text = t

	s.draw()
}

func (s *status) clear() {
	s.set("")
}

func (s *status) draw() {
	s.textView.Clear()
	fmt.Fprintf(s.textView, "@%s %s", shared.api.CurrentUser.UserName, s.text)
}
