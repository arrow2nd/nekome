package app

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalOpt : モーダルの設定
type ModalOpt struct {
	title  string
	onDone func()
}

type view struct {
	pageView *tview.Pages
	tabView  *tview.TextView
	modal    *tview.Modal
	pages    map[string]page
	tabs     []string
	mu       sync.Mutex
}

func newView() *view {
	v := &view{
		pageView: tview.NewPages(),
		tabView:  tview.NewTextView(),
		modal:    tview.NewModal(),
		pages:    map[string]page{},
		tabs:     []string{},
	}

	v.tabView.
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignLeft).
		SetHighlightedFunc(v.handleTabHighlight).
		SetBackgroundColor(tcell.ColorDefault)

	v.modal.
		AddButtons([]string{"No", "Yes"}).
		SetTextColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(v.handleModalKeyEvent)

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
		fmt.Fprintf(v.tabView, `["%s"] %s [""]`, createPageTag(i), name)

		if i < len(v.tabs)-1 {
			fmt.Fprint(v.tabView, "|")
		}
	}
}

// SetInputCapture : キーイベントハンドラを設定
func (v *view) SetInputCapture(f func(*tcell.EventKey) *tcell.EventKey) {
	v.pageView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// モーダル表示中は操作を受け付けない
		if v.modal.HasFocus() {
			return event
		}

		return f(event)
	})
}

// AddPage : ページを追加
func (v *view) AddPage(p page, focus bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	// タブを追加
	v.tabs = append(v.tabs, p.GetName())
	v.drawTab()

	pageID := createPageTag(len(v.tabs) - 1)

	// ページを追加
	v.pages[pageID] = p
	v.pageView.AddPage(pageID, p.GetPrimivite(), true, focus)

	if focus {
		v.tabView.Highlight(pageID)
	}

	go p.Load()
}

// RemovePage : ページを削除
func (v *view) RemovePage(name string) {
	v.pageView.RemovePage(name)
}

// PopupModal : モーダルを表示
func (v *view) PopupModal(o *ModalOpt) {
	v.modal.
		SetText(o.title).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				o.onDone()
			}
			v.RemovePage("modal")
		})

	v.pageView.AddPage("modal", v.modal, true, true)
}

// selectPrevTab : 前のタブを選択
func (v *view) selectPrevTab() {
	index := getHighlightId(v.tabView.GetHighlights())
	if index == -1 {
		return
	}

	pageCount := v.pageView.GetPageCount()

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

	pageCount := v.pageView.GetPageCount()

	index = (index + 1) % pageCount

	v.tabView.Highlight(createPageTag(index))
}

// handleTabHighlight : タブがハイライトされたときのコールバック
func (v *view) handleTabHighlight(added, removed, remaining []string) {
	// ハイライトされたタブまでスクロール
	v.tabView.ScrollToHighlight()

	// ページを切り替え
	v.pageView.SwitchToPage(added[0])
	v.pages[added[0]].OnVisible()
}

// handleModalKeyEvent : モーダルのキーイベントハンドラ
func (v *view) handleModalKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	keyRune := event.Rune()

	// 左キーの入力イベントに置換
	if keyRune == 'h' || keyRune == 'j' {
		return tcell.NewEventKey(tcell.KeyLeft, keyRune, tcell.ModNone)
	}

	// 右キーの入力イベントに置換
	if keyRune == 'l' || keyRune == 'k' {
		return tcell.NewEventKey(tcell.KeyRight, keyRune, tcell.ModNone)
	}

	return event
}
