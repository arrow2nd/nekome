package ui

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type view struct {
	pages   *tview.Pages
	tabView *tview.TextView
	tabs    []string
	mu      sync.Mutex
}

func newView() *view {
	v := &view{
		pages:   tview.NewPages(),
		tabView: tview.NewTextView(),
		tabs:    []string{},
	}

	v.tabView.
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignLeft).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			v.pages.SwitchToPage(added[0])
		}).
		SetBackgroundColor(tcell.ColorDefault)

	return v
}

// createPageTag : ページ管理用のタグ文字列を作成
func createPageTag(id int) string {
	return fmt.Sprintf("page_%d", id)
}

// drawTab : タブを描画
func (v *view) drawTab() {
	v.tabView.Clear()

	for i, name := range v.tabs {
		fmt.Fprintf(v.tabView, `["%s"] %s `, createPageTag(i), name)
	}
}

// AddPage : ページを追加
func (v *view) AddPage(p page, focus bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	// タブを追加
	v.tabs = append(v.tabs, p.GetName())
	v.drawTab()

	// ページを追加
	pageID := v.createPageTag(len(v.tabs) - 1)
	v.pages.AddPage(pageID, p.GetPrimivite(), true, focus)

	if focus {
		v.tabView.Highlight(pageID)
	}

	go p.Load()
}

// selectPrevTab : 前のタブを選択
func (v *view) selectPrevTab() {
	index := getHighlightId(v.tabView.GetHighlights())
	if index == -1 {
		return
	}

	pageCount := v.pages.GetPageCount()

	if index--; index < 0 {
		index = pageCount - 1
	}

	v.tabView.Highlight(v.createPageTag(index))
}

// selectNextTab : 次のタブを選択
func (v *view) selectNextTab() {
	index := getHighlightId(v.tabView.GetHighlights())
	if index == -1 {
		return
	}

	pageCount := v.pages.GetPageCount()

	index = (index + 1) % pageCount

	v.tabView.Highlight(v.createPageTag(index))
}
