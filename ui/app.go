package ui

import (
	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI ユーザインターフェース
type UI struct {
	app         *tview.Application
	view        *view
	statusBar   *statusBar
	commandLine *tview.InputField
}

// New 生成
func New() *UI {
	return &UI{
		app:         tview.NewApplication(),
		view:        newView(),
		statusBar:   newStatusBar(),
		commandLine: tview.NewInputField(),
	}
}

// Init 初期化
func (u *UI) Init(a *api.API, c *config.Config) {
	// 共有
	shared.api = a
	shared.conf = c

	// 配色設定
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault

	// ページ
	home := newTimelinePage(HomeTL)
	mention := newTimelinePage(MentionTL)

	u.view.addPage("Home", home, true)
	u.view.addPage("Mention", mention, false)

	u.view.pages.SetInputCapture(u.handlePageKeyEvent)

	// ステータスバー
	u.statusBar.draw()

	// 入力フィールド
	u.initCommandLine()

	// 画面レイアウト
	// NOTE: 追加順がキーハンドラの優先順になるっぽい
	layout := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetBorders(false).
		AddItem(u.view.tabTextView, 0, 0, 1, 1, 0, 0, false).
		AddItem(u.statusBar.flex, 2, 0, 1, 1, 0, 0, false).
		AddItem(u.commandLine, 3, 0, 1, 1, 0, 0, false).
		AddItem(u.view.pages, 1, 0, 1, 1, 0, 0, true)

	u.app.SetRoot(layout, true)
	u.app.SetInputCapture(u.handleGlobalKeyEvents)
}

// Run 実行
func (u *UI) Run() error {
	go u.eventReciever()
	return u.app.Run()
}

func (u *UI) eventReciever() {
	for {
		select {
		case status := <-shared.stateCh:
			u.commandLine.SetPlaceholder(status)
			u.app.Draw()
		case <-shared.appDrawCh:
			u.app.Draw()
		}
	}
}

func (u *UI) redraw() {
	// NOTE: 絵文字の表示幅問題で表示が崩れてしまう問題への暫定的な対応
	// https://github.com/rivo/tview/issues/693

	pageId, _ := u.view.pages.GetFrontPage()
	if pageId == "" {
		shared.setStatus("No page to redraw")
		return
	}

	// NOTE: 再描画する方法が見当たらなかったので、該当ページを非表示にして再度表示
	// することで実質的に再描画を行う
	u.view.pages.HidePage(pageId)
	u.app.ForceDraw()
	u.view.pages.ShowPage(pageId)

	shared.setStatus("Redraw!")
}

func (u *UI) handleGlobalKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyCtrlC {
		u.app.Stop()
		return nil
	}

	return event
}

func (u *UI) handlePageKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyLeft:
		u.view.selectPrevTab()
		return nil
	case tcell.KeyRight:
		u.view.selectNextTab()
		return nil
	case tcell.KeyCtrlL:
		u.redraw()
		return nil
	case tcell.KeyRune:
		switch event.Rune() {
		case 'h':
			u.view.selectPrevTab()
			return nil
		case 'l':
			u.view.selectNextTab()
			return nil
		case ':':
			u.app.SetFocus(u.commandLine)
			return nil
		}
	}

	return event
}
