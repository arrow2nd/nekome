package app

import (
	"strings"

	"github.com/rivo/tview"
)

type helpPage struct {
	*basePage
	textView *tview.TextView
}

// newHelpPage : ヘルプページ生成
func newHelpPage(name, text string) *helpPage {
	tabName := shared.conf.Settings.Texts.TabHelp
	tabName = strings.Replace(tabName, "{name}", name, 1)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetText(text)

	p := &helpPage{
		basePage: newBasePage(tabName),
		textView: textView,
	}

	p.SetFrame(p.textView)

	return p
}

// Load : pageのインターフェースを満たすためのダミー
func (h *helpPage) Load() {}
