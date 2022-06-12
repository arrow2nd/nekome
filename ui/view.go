package ui

import (
	"fmt"
	"sync"

	"github.com/rivo/tview"
)

// view ページ・タブ管理
type view struct {
	pages       *tview.Pages
	tabTextView *tview.TextView
	tabs        []string
	mu          sync.Mutex
}

func newView() *view {
	v := &view{
		pages:       tview.NewPages(),
		tabTextView: tview.NewTextView(),
		tabs:        []string{},
	}

	v.tabTextView.
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignLeft).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			v.pages.SwitchToPage(added[0])
		})

	return v
}

func (v *view) drawTab() {
	v.tabTextView.Clear()

	for i, name := range v.tabs {
		fmt.Fprintf(v.tabTextView, `["%s"] %s `, v.createPageId(i), name)
	}
}

func (v *view) createPageId(id int) string {
	return fmt.Sprintf("page_%d", id)
}

func (v *view) addTab(tabName string, focus bool) string {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.tabs = append(v.tabs, tabName)

	newPageId := v.createPageId(len(v.tabs) - 1)

	if focus {
		v.tabTextView.Highlight(newPageId)
	}

	return newPageId
}

func (v *view) addPage(tabName string, item tview.Primitive, focus bool) {
	pageId := v.addTab(tabName, focus)
	v.pages.AddPage(pageId, item, true, focus)

	v.drawTab()
}

func (v *view) selectPrevTab() {
	index := getHighlightId(v.tabTextView.GetHighlights())
	if index == -1 {
		return
	}

	pageCount := v.pages.GetPageCount()

	if index--; index < 0 {
		index = pageCount - 1
	}

	v.tabTextView.Highlight(v.createPageId(index))
}

func (v *view) selectNextTab() {
	index := getHighlightId(v.tabTextView.GetHighlights())
	if index == -1 {
		return
	}

	pageCount := v.pages.GetPageCount()

	index = (index + 1) % pageCount

	v.tabTextView.Highlight(v.createPageId(index))
}
