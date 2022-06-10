package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type statusBar struct {
	textView *tview.TextView
}

func newStatusBar() *statusBar {
	s := &statusBar{
		textView: tview.NewTextView(),
	}

	s.textView.SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft).
		SetBackgroundColor(tcell.ColorDarkCyan)

	return s
}

func (s *statusBar) draw() {
	s.textView.Clear()
	fmt.Fprintf(s.textView, "@%s", shared.api.CurrentUser.UserName)
}
