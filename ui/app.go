package ui

import (
	"sync"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

// UI :  ユーザインターフェース
type UI struct {
	app         *tview.Application
	view        *view
	statusBar   *statusBar
	commandLine *tview.InputField
	mu          sync.Mutex
}

// New : 生成
func New() *UI {
	return &UI{
		app:         tview.NewApplication(),
		view:        newView(),
		statusBar:   newStatusBar(),
		commandLine: tview.NewInputField(),
	}
}

// Init : 初期化
func (u *UI) Init(a *api.API, c *config.Config) {
	// 日本語環境等での罫線の乱れ対策
	runewidth.DefaultCondition.EastAsianWidth = false

	// 共有
	shared.api = a
	shared.conf = c

	// 配色設定
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault

	// ページ
	home := newTimelinePage(homeTL)
	mention := newTimelinePage(mentionTL)
	user := newUserPage("imas_official")
	userB := newUserPage("arrow_2nd")

	u.view.AddPage(home, true)
	u.view.AddPage(mention, false)
	u.view.AddPage(user, false)
	u.view.AddPage(userB, false)

	u.view.pages.SetInputCapture(u.handlePageKeyEvent)

	// ステータスバー
	u.statusBar.Draw()

	// コマンドライン
	u.initCommandLine()

	// 画面レイアウト
	// NOTE: 追加順がキーハンドラの優先順になるっぽい
	layout := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetBorders(false).
		AddItem(u.view.tabView, 0, 0, 1, 1, 0, 0, false).
		AddItem(u.statusBar.flex, 2, 0, 1, 1, 0, 0, false).
		AddItem(u.commandLine, 3, 0, 1, 1, 0, 0, false).
		AddItem(u.view.pages, 1, 0, 1, 1, 0, 0, true)

	u.app.SetRoot(layout, true)
	u.app.SetInputCapture(u.handleGlobalKeyEvents)
}

// Run : 実行
func (u *UI) Run() error {
	go u.eventReciever()
	return u.app.Run()
}

// eventReciever : イベントレシーバ
func (u *UI) eventReciever() {
	for {
		select {
		case status := <-shared.chStatus:
			u.setStatusMessage(status)
		}
	}
}

// redraw : 全体を再描画
func (u *UI) redraw() {
	// NOTE: 絵文字の表示幅問題で表示が崩れてしまう問題への暫定的な対応
	// https://github.com/rivo/tview/issues/693

	pageId, _ := u.view.pages.GetFrontPage()
	if pageId == "" {
		shared.SetErrorStatus("App", "no page to redraw")
		return
	}

	// 一度非表示にして画面をクリア
	u.view.pages.HidePage(pageId)

	// 強制的に再描画して画面を再表示
	u.app.ForceDraw()
	u.view.pages.ShowPage(pageId)
}

// handleGlobalKeyEvents : アプリ全体のキーハンドラ
func (u *UI) handleGlobalKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	// アプリを終了
	if key == tcell.KeyCtrlC || key == tcell.KeyCtrlQ {
		u.app.Stop()
		return nil
	}

	return event
}

// handlePageKeyEvent : ページビューのキーハンドラ
func (u *UI) handlePageKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	keyRune := event.Rune()

	// 左のタブを選択
	if key == tcell.KeyLeft || keyRune == 'h' {
		u.view.selectPrevTab()
		return nil
	}

	// 右のタブを選択
	if key == tcell.KeyRight || keyRune == 'l' {
		u.view.selectNextTab()
		return nil
	}

	// 再描画
	if key == tcell.KeyCtrlL {
		u.redraw()
		return nil
	}

	// コマンドラインへフォーカスを移動
	if keyRune == ':' {
		u.app.SetFocus(u.commandLine)
		return nil
	}

	return event
}
