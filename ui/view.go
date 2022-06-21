package ui

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type view struct {
	pagesView *tview.Pages
	tabView   *tview.TextView
	pages     map[string]page
	tabNames  []string
	mu        sync.Mutex
}

func newView() *view {
	v := &view{
		pagesView: tview.NewPages(),
		tabView:   tview.NewTextView(),
		pages:     map[string]page{},
		tabNames:  []string{},
	}

	v.tabView.
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignLeft).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			v.pagesView.SwitchToPage(added[0])
			v.pages[added[0]].OnVisible()
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

	for i, name := range v.tabNames {
		fmt.Fprintf(v.tabView, `["%s"] %s `, createPageTag(i), name)
	}
}

// AddPage : ページを追加
func (v *view) AddPage(p page, focus bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	// タブを追加
	v.tabNames = append(v.tabNames, p.GetName())
	v.drawTab()

	pageID := createPageTag(len(v.tabNames) - 1)

	// ページを追加
	v.pages[pageID] = p
	v.pagesView.AddPage(pageID, p.GetPrimivite(), true, focus)

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

	pageCount := v.pagesView.GetPageCount()

	if index--; index < 0 {
		index = pageCount - 1
	}

	v.tabView.Highlight(createPageTag(index))
}

// selectNextTab : 次のタブを選択
func (v *view) selectNextTab() {
	index := getHighlightId(v.tabView.GetHighlights())
	if index == -1 {
		return
	}

	pageCount := v.pagesView.GetPageCount()

	index = (index + 1) % pageCount

	v.tabView.Highlight(createPageTag(index))
}
