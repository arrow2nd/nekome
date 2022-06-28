package app

import (
	"strings"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

// App : アプリケーション
type App struct {
	app         *tview.Application
	view        *view
	statusBar   *statusBar
	commandLine *commandLine
}

// New : 生成
func New() *App {
	return &App{
		app:         tview.NewApplication(),
		view:        newView(),
		statusBar:   newStatusBar(),
		commandLine: newCommandLine(),
	}
}

// Init : 初期化
func (a *App) Init(app *api.API, conf *config.Config) {
	// 日本語環境等での罫線の乱れ対策
	runewidth.DefaultCondition.EastAsianWidth = false

	// 全体共有
	shared.api = app
	shared.conf = conf

	// 配色設定
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault

	// ページビュー
	a.view.SetInputCapture(a.handlePageKeyEvent)

	// ステータスバー
	a.statusBar.DrawAccountInfo()

	// コマンドライン
	go func() {
		if err := a.commandLine.SetListCompleteItems(); err != nil {
			shared.SetErrorStatus("Init - CommandLine", err.Error())
		}
	}()

	// 画面レイアウト
	// NOTE: 追加順がキーハンドラの優先順になるっぽい
	layout := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetBorders(false).
		AddItem(a.view.tabView, 0, 0, 1, 1, 0, 0, false).
		AddItem(a.statusBar.flex, 2, 0, 1, 1, 0, 0, false).
		AddItem(a.commandLine.inputField, 3, 0, 1, 1, 0, 0, false).
		AddItem(a.view.pageView, 1, 0, 1, 1, 0, 0, true)

	a.app.SetRoot(layout, true).
		SetInputCapture(a.handleGlobalKeyEvents)

	// コマンドを実行
	for _, c := range shared.conf.Settings.Feature.RunCommands {
		if err := a.ExecCommand(strings.Split(c, " ")); err != nil {
			shared.SetErrorStatus("Command", err.Error())
		}
	}
}

// Run : 実行
func (a *App) Run() error {
	go a.eventReciever()
	return a.app.Run()
}

// eventReciever : イベントレシーバ
func (a *App) eventReciever() {
	for {
		select {
		case status := <-shared.chStatus:
			a.commandLine.UpdateStatusMessage(status)
			a.app.Draw()
		case indicator := <-shared.chIndicator:
			a.statusBar.DrawPageIndicator(indicator)
			a.app.Draw()
		case opt := <-shared.chPopupModal:
			a.view.PopupModal(opt)
			a.app.Draw()
		case cmd := <-shared.chExecCommand:
			if err := a.ExecCommand(strings.Split(cmd, " ")); err != nil {
				shared.SetErrorStatus("Command", err.Error())
			}
		case cmd := <-shared.chInputCommand:
			a.app.SetFocus(a.commandLine.inputField)
			a.commandLine.SetText(cmd)
			a.app.Draw()
		case <-shared.chFocusPageView:
			a.app.SetFocus(a.view.pageView)
			a.app.Draw()
		}
	}
}

// redraw : アプリ全体を再描画
func (a *App) redraw() {
	// NOTE: 絵文字の表示幅問題で表示が崩れてしまう問題への暫定的な対応
	// https://github.com/rivo/tview/issues/693

	pageId, _ := a.view.pageView.GetFrontPage()
	if pageId == "" {
		shared.SetErrorStatus("App", "no page to redraw")
		return
	}

	// 一度非表示にして画面をクリア
	a.view.pageView.HidePage(pageId)

	// 強制的に再描画して画面を再表示
	a.app.ForceDraw()
	a.view.pageView.ShowPage(pageId)
}

// handleGlobalKeyEvents : アプリ全体のキーハンドラ
func (a *App) handleGlobalKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()

	// アプリを終了
	if key == tcell.KeyCtrlQ {
		a.commandLine.Blur()
		a.quitApp()
		return nil
	}

	return event
}

// handlePageKeyEvent : ページビューのキーハンドラ
func (a *App) handlePageKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	keyRune := event.Rune()

	// 左のタブを選択
	if key == tcell.KeyLeft || keyRune == 'h' {
		a.view.selectPrevTab()
		return nil
	}

	// 右のタブを選択
	if key == tcell.KeyRight || keyRune == 'l' {
		a.view.selectNextTab()
		return nil
	}

	// 再描画
	if key == tcell.KeyCtrlL {
		a.redraw()
		return nil
	}

	// コマンドラインへフォーカスを移動
	if keyRune == ':' {
		a.app.SetFocus(a.commandLine.inputField)
		return nil
	}

	// ページを削除
	if key == tcell.KeyCtrlW {
		a.view.RemoveCurrentPage()
		return nil
	}

	return event
}
