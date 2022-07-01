package app

import (
	"os"
	"strings"

	"github.com/arrow2nd/nekome/api"
	"github.com/arrow2nd/nekome/config"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// App : アプリケーション
type App struct {
	app         *tview.Application
	cmd         *cobra.Command
	view        *view
	statusBar   *statusBar
	commandLine *commandLine
}

// New : 生成
func New() *App {
	return &App{
		app:         tview.NewApplication(),
		cmd:         newCmd(),
		view:        newView(),
		statusBar:   newStatusBar(),
		commandLine: newCommandLine(),
	}
}

// Init : 初期化
func (a *App) Init(app *api.API, conf *config.Config) {
	// 全体共有
	shared.isCommandLineMode = len(os.Args[1:]) > 0
	shared.api = app
	shared.conf = conf

	// コマンド
	a.initCmd()

	// コマンドラインモードならUIの初期化をスキップ
	if shared.isCommandLineMode {
		return
	}

	// 日本語環境等での罫線の乱れ対策
	// https://github.com/mattn/go-runewidth/issues/14
	runewidth.DefaultCondition.EastAsianWidth = !shared.conf.Settings.Feature.IsLocaleCJK

	// 背景色
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault

	// ページビュー
	a.view.SetInputCapture(a.handlePageKeyEvent)

	// ステータスバー
	a.statusBar.Init()
	a.statusBar.DrawAccountInfo()

	// コマンドライン
	a.commandLine.Init()
	go func() {
		cmds := a.cmd.Commands()
		if err := a.commandLine.SetListCompleteItems(cmds); err != nil {
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

	a.app.
		SetRoot(layout, true).
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
	// コマンドラインモード
	if shared.isCommandLineMode {
		return a.ExecCommand(os.Args[1:])
	}

	go a.eventReciever()
	return a.app.Run()
}

// ExecCommand : コマンドを実行
func (a *App) ExecCommand(args []string) error {
	a.cmd.SetArgs(args)
	return a.cmd.Execute()
}

// quitApp : アプリを終了
func (a *App) quitApp() {
	// 確認画面が不要ならそのまま終了
	if !shared.conf.Settings.Feature.Confirm["Quit"] {
		a.app.Stop()
		return
	}

	a.view.PopupModal(&ModalOpt{
		title:  "Do you want to quit the app?",
		onDone: a.app.Stop,
	})
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

// handleGlobalKeyEvents : アプリ全体のキーハンドラ
func (a *App) handleGlobalKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	keyRune := event.Rune()

	// アプリを終了
	if key == tcell.KeyCtrlQ {
		a.commandLine.Blur()
		a.quitApp()
		return nil
	}

	// ショートカットのヘルプ
	if keyRune == '?' {
		shared.RequestExecCommand("helpshortcuts")
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
		a.app.Sync()
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
