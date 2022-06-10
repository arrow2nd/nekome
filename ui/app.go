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
	home := newHomeTimeline()
	home.init()

	u.view.addPage("Home", home.frame, true)
	u.view.addPage("Mention", home.frame, false)

	u.view.pages.SetInputCapture(u.handlePageKeyEvent)

	// ステータスバー
	u.statusBar.draw()

	// 入力フィールド
	u.initCommandLine()

	// 画面レイアウト

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(u.view.tabTextView, 2, 1, false).
		AddItem(u.view.pages, 0, 1, true).
		AddItem(u.statusBar.textView, 1, 1, false).
		AddItem(u.commandLine, 1, 1, false)

	u.app.SetRoot(layout, true)

	// マウス操作有効化
	u.app.EnableMouse(true).
		SetBeforeDrawFunc(func(screen tcell.Screen) bool {
			screen.Clear()
			return false
		})
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
		}
	}
}

func (u *UI) handlePageKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyLeft:
		u.view.selectPrevTab()
		return nil
	case tcell.KeyRight:
		u.view.selectNextTab()
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
