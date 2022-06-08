package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

type tabBar struct {
	textView *tview.TextView
	tabs     []string
}

func newTabBar() *tabBar {
	tb := &tabBar{
		textView: tview.NewTextView(),
		tabs:     []string{},
	}

	tb.textView.
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignLeft)

	return tb
}

func (tb *tabBar) SetTab(tabname []string) {
	for i, name := range tabname {
		fmt.Fprintf(tb.textView, `["page_%d"] %s `, i, name)
	}

	tb.textView.Highlight("page_0")
}
