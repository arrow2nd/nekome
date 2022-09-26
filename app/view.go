package app

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// tabMove : タブの移動量
type tabMove int

const (
	TabMovePrev tabMove = -1
	TabMoveNext tabMove = 1
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
	mainView *tview.Flex
	pageView *tview.Pages
	tabView  *tview.TextView
	textArea *tview.TextArea
	modal    *tview.Modal
	pages    map[string]page
	tabs     []*tab
	tabIndex int
	mu       sync.Mutex
}

func newView() *view {
	v := &view{
		mainView: tview.NewFlex(),
		pageView: tview.NewPages(),
		tabView:  tview.NewTextView(),
		textArea: tview.NewTextArea(),
		modal:    tview.NewModal(),
		pages:    map[string]page{},
		tabs:     []*tab{},
		tabIndex: 0,
	}

	v.mainView.
		SetDirection(tview.FlexRow).
		AddItem(v.pageView, 0, 1, true).
		AddItem(v.textArea, 0, 0, false)

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
		SetInputCapture(v.handleModalKeyEvents)

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

// SetInputCapture : キーハンドラを設定
func (v *view) SetInputCapture(f func(*tcell.EventKey) *tcell.EventKey) {
	v.mainView.SetInputCapture(f)
}

// AddPage : ページを追加
func (v *view) AddPage(p page, focus bool) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	newTab := &tab{
		id:   getMD5(p.GetName()),
		name: p.GetName(),
	}

	// ページが重複する場合、既にあるページに移動
	if _, ok := v.pages[newTab.id]; ok {
		tabIndex, found := find(v.tabs, func(e *tab) bool { return e.id == newTab.id })
		if !found {
			return fmt.Errorf("Failed to add page (%s)", newTab.name)
		}

		v.tabView.Highlight(newTab.id)
		v.tabIndex = tabIndex

		return nil
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

	// 前のタブを選択
	if v.tabIndex--; v.tabIndex < 0 {
		v.tabIndex = 0
	}

	v.tabView.Highlight(v.tabs[v.tabIndex].id)
}

// MoveTab : タブを移動する
func (v *view) MoveTab(move tabMove) {
	prevTabIndex := v.tabIndex
	v.tabIndex += int(move)

	// 範囲内に丸める
	if max := v.pageView.GetPageCount(); v.tabIndex < 0 {
		v.tabIndex = max - 1
	} else if v.tabIndex >= max {
		v.tabIndex = 0
	}

	// 移動前と同じなら中断
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
			shared.SetDisablePageKeyEvent(false)
		}).
		SetButtonBackgroundColor(tcell.ColorDefault)

	v.pageView.AddPage("modal", v.modal, true, true)

	shared.SetDisablePageKeyEvent(true)
}

// handleModalKeyEvents : モーダルのキーハンドラ
func (v *view) handleModalKeyEvents(event *tcell.EventKey) *tcell.EventKey {
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

// ShowTextArea : テキストエリアを表示
func (v *view) ShowTextArea(title string, onSubmit func(s string)) {
	v.textArea.
		SetText("", false).
		SetTextStyle(tcell.StyleDefault).
		SetTitle(fmt.Sprintf(" %s (Press ESC to close, press Ctrl-P to post) ", title)).
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1).
		SetBorder(true).
		SetTitleColor(tcell.ColorDefault).
		SetBorderColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			key := event.Key()

			// 閉じる
			if key == tcell.KeyEsc {
				v.HiddenTextArea()
				return nil
			}

			// 送信
			if key == tcell.KeyCtrlP {
				v.HiddenTextArea()
				onSubmit(v.textArea.GetText())
				return nil
			}

			return event
		})

	v.mainView.ResizeItem(v.textArea, 0, 1)

	shared.RequestFocusPrimitive(v.textArea)
	shared.SetDisablePageKeyEvent(true)
}

// HiddenTextArea : テキストエリアを非表示
func (v *view) HiddenTextArea() {
	v.mainView.ResizeItem(v.textArea, 0, 0)

	shared.RequestFocusPrimitive(v.pageView)
	shared.SetDisablePageKeyEvent(false)
}
