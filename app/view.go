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
	text   string
	onDone func()
}

// tab : タブ
type tab struct {
	id   string
	name string
}

// view : ページの表示域
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

	t := shared.conf.Style.App.Tab

	for i, tab := range v.tabs {
		fmt.Fprintf(v.tabView, `[%s]["%s"] %s [""][-:-:-]`, t, tab.id, tab.name)

		// タブが2個以上あるならセパレータを挿入
		if i < len(v.tabs)-1 {
			fmt.Fprint(v.tabView, shared.conf.Settings.Appearance.TabSeparate)
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

	// フォーカスが当たっているならタブをハイライト
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

// Reset : リセット
func (v *view) Reset() {
	// ページを削除
	for id := range v.pages {
		v.pageView.RemovePage(id)
	}
	v.pages = map[string]page{}

	// タブを削除
	v.tabs = []*tab{}
	v.tabView.SetText("")
	v.tabIndex = 0
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
	v.pages[id].OnDelete()
	delete(v.pages, id)

	// 1つ前のタブを選択
	if v.tabIndex--; v.tabIndex < 0 {
		v.tabIndex = 0
	}

	v.tabView.Highlight(v.tabs[v.tabIndex].id)
}

// PopupModal : モーダルを表示
func (v *view) PopupModal(o *ModalOpt) {
	message := o.title

	if o.text != "" {
		hr := createSeparator("-", 4)
		message = fmt.Sprintf("%s\n%s\n%s", o.title, hr, o.text)
	}

	v.modal.
		SetFocus(0).
		SetText(message).
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
	prevTabIndex := v.tabIndex
	pageCount := v.pageView.GetPageCount()

	if v.tabIndex--; v.tabIndex < 0 {
		v.tabIndex = pageCount - 1
	}

	if v.tabIndex == prevTabIndex {
		return
	}

	v.tabView.Highlight(v.tabs[v.tabIndex].id)
}

// selectNextTab : 次のタブを選択
func (v *view) selectNextTab() {
	prevTabIndex := v.tabIndex
	pageCount := v.pageView.GetPageCount()

	v.tabIndex = (v.tabIndex + 1) % pageCount

	if v.tabIndex == prevTabIndex {
		return
	}

	v.tabView.Highlight(v.tabs[v.tabIndex].id)
}

// handleTabHighlight : タブがハイライトされたときのコールバック
func (v *view) handleTabHighlight(added, removed, remaining []string) {
	// ハイライトされたタブまでスクロール
	v.tabView.ScrollToHighlight()

	// 前のページを非アクティブにする
	if len(removed) > 0 {
		if page, ok := v.pages[removed[0]]; ok {
			page.OnInactive()
		}
	}

	// ページを切り替え
	v.pageView.SwitchToPage(added[0])
	v.pages[added[0]].OnActive()
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
