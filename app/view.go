package app

import (
	"errors"
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

type tab struct {
	id   string
	name string
}

type view struct {
	pageView *tview.Pages
	tabView  *tview.TextView
	modal    *tview.Modal
	pages    map[string]page
	tabs     []*tab
	tabIndex int
	mu       sync.Mutex
}

func newView() *view {
	v := &view{
		pageView: tview.NewPages(),
		tabView:  tview.NewTextView(),
		modal:    tview.NewModal(),
		pages:    map[string]page{},
		tabs:     []*tab{},
		tabIndex: 0,
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

	for i, tab := range v.tabs {
		fmt.Fprintf(v.tabView, `["%s"] %s [""]`, tab.id, tab.name)

		if i < len(v.tabs)-1 {
			fmt.Fprint(v.tabView, shared.conf.Settings.Apperance.TabSeparate)
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
func (v *view) AddPage(p page, focus bool) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	newTab := &tab{
		id:   getMD5(p.GetName()),
		name: p.GetName(),
	}

	// ページの重複チェック
	if _, ok := v.pages[newTab.id]; ok {
		return errors.New("this page already exists")
	}

	// ページを追加
	v.pages[newTab.id] = p
	v.pageView.AddPage(newTab.id, p.GetPrimivite(), true, focus)

	if focus {
		v.tabView.Highlight(newTab.id)
		v.tabIndex = v.pageView.GetPageCount() - 1
	}

	// タブを追加
	v.tabs = append(v.tabs, newTab)
	v.drawTab()

	go p.Load()

	return nil
}

// RemoveCurrentPage : 現在のページを削除
func (v *view) RemoveCurrentPage() {
	// ページが1つのみなら削除しない
	if v.pageView.GetPageCount() == 1 {
		shared.SetErrorStatus("App", "last page cannot be deleted")
		return
	}

	id, _ := v.pageView.GetFrontPage()
	name := v.pages[id].GetName()

	// タブを削除
	newTabs := []*tab{}

	for _, tab := range v.tabs {
		if tab.name != name {
			newTabs = append(newTabs, tab)
		}
	}

	v.tabs = newTabs

	// 再描画して反映
	v.drawTab()

	// ページを削除
	v.pageView.RemovePage(id)
	delete(v.pages, id)

	// 1つ前のタブを選択
	v.tabIndex--
	if v.tabIndex < 0 {
		v.tabIndex = 0
	}

	v.tabView.Highlight(v.tabs[v.tabIndex].id)
}

// PopupModal : モーダルを表示
func (v *view) PopupModal(o *ModalOpt) {
	v.modal.
		SetFocus(0).
		SetText(o.title).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				o.onDone()
			}
			v.pageView.RemovePage("modal")
		})

	v.pageView.AddPage("modal", v.modal, true, true)
}

// selectPrevTab : 前のタブを選択
func (v *view) selectPrevTab() {
	pageCount := v.pageView.GetPageCount()

	if v.tabIndex--; v.tabIndex < 0 {
		v.tabIndex = pageCount - 1
	}

	v.tabView.Highlight(v.tabs[v.tabIndex].id)
}

// selectNextTab : 次のタブを選択
func (v *view) selectNextTab() {
	pageCount := v.pageView.GetPageCount()

	v.tabIndex = (v.tabIndex + 1) % pageCount

	v.tabView.Highlight(v.tabs[v.tabIndex].id)
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

	// hjを左キーの入力イベントに置換
	if keyRune == 'h' || keyRune == 'j' {
		return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
	}

	// klを右キーの入力イベントに置換
	if keyRune == 'k' || keyRune == 'l' {
		return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
	}

	return event
}
