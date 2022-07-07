package app

import (
	"strings"

	"github.com/rivo/tview"
)

type docsPage struct {
	*basePage
	textView *tview.TextView
}

func newDocsPage(name, text string) *docsPage {
	tabName := shared.conf.Settings.Texts.TabDocs
	tabName = strings.Replace(tabName, "{name}", name, 1)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetText(text)

	p := &docsPage{
		basePage: newBasePage(tabName),
		textView: textView,
	}

	p.SetFrame(p.textView)

	return p
}

// Load : pageのインターフェースを満たすためのダミー
func (d *docsPage) Load() {}
