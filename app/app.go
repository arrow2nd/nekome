package app

import (
	"fmt"
	"os"

	"github.com/arrow2nd/nekome/cli"
	"github.com/arrow2nd/nekome/config"
	"github.com/arrow2nd/nekome/log"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

var version = "develop"

// App : アプリケーション
type App struct {
	app         *tview.Application
	cmd         *cli.Command
	view        *view
	statusBar   *statusBar
	commandLine *commandLine
	args        []string
}

// New : 生成
func New() *App {
	return &App{
		app:         tview.NewApplication(),
		cmd:         newCmd(),
		view:        newView(),
		statusBar:   newStatusBar(),
		commandLine: newCommandLine(),
		args:        []string{},
	}
}

// Init : 初期化
func (a *App) Init() error {
	// 設定読み込み
	if err := a.loadConfig(); err != nil {
		return err
	}

	// フラグをパースして対応する処理を実行
	if err := a.parseRuntimeFlags(); err != nil {
		return err
	}

	// コマンド初期化
	a.initCmd()

	// コマンドラインモードならUIの初期化をスキップ
	if shared.isCommandLineMode {
		return nil
	}

	// 日本語環境等での罫線の乱れ対策
	// https://github.com/mattn/go-runewidth/issues/14
	runewidth.DefaultCondition.EastAsianWidth = !shared.conf.Settings.Feature.IsLocaleCJK

	// 背景色
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault

	// ページのキーハンドラを設定
	a.view.SetInputCapture(a.handlePageKeyEvent)

	// ステータスバー初期化
	a.statusBar.Init()
	a.statusBar.DrawAccountInfo()

	// コマンドライン初期化
	a.commandLine.Init()
	a.initAutocomplate()

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

	a.runStartupCommands()

	return nil
}

// loadConfig : 設定を読み込む
func (a *App) loadConfig() error {
	shared.conf = config.New()

	// 環境設定
	if err := shared.conf.LoadSettings(); err != nil {
		return err
	}

	// スタイル
	if err := shared.conf.LoadStyle(); err != nil {
		return err
	}

	// 認証情報
	ok, err := shared.conf.LoadCred()
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	// 認証情報が無い場合、新規追加
	return addAccount(true)
}

// parseRuntimeFlags : 実行時のフラグをパース
func (a *App) parseRuntimeFlags() error {
	// フラグをパース
	f := a.cmd.NewFlagSet()
	if err := f.Parse(os.Args[1:]); err != nil {
		return err
	}

	// コマンドラインモードかどうか
	shared.isCommandLineMode = f.NArg() > 0 || f.Changed("help")

	// ヘルプフラグが指定されているならログインは行わない
	if f.Changed("help") {
		a.args = os.Args[1:]
		return nil
	}

	// バージョンフラグが指定されているなら表示して終了
	if f.Changed("version") {
		log.LogExit(fmt.Sprintf("🐈 nekome v.%s", version))
	}

	// 引数を保存
	a.args = f.Args()

	// ログイン処理
	user, _ := f.GetString("user")
	return loginAccount(user)
}

// initAutocomplate : 入力補完を初期化
func (a *App) initAutocomplate() {
	cmds := a.cmd.GetChildrenNames(true)
	if err := a.commandLine.SetAutocompleteItems(cmds); err != nil {
		shared.SetErrorStatus("Init - CommandLine", err.Error())
	}

}

// runStartupCommands : 起動時に実行するコマンドを実行
func (a *App) runStartupCommands() {
	for _, c := range shared.conf.Settings.Feature.Startup {
		if err := a.RunCommand(c); err != nil {
			shared.SetErrorStatus("Command", err.Error())
		}
	}
}

// Run : 実行
func (a *App) Run() error {
	// コマンドラインモード
	if shared.isCommandLineMode {
		return a.cmd.Execute(a.args)
	}

	go a.eventReciever()
	return a.app.Run()
}

// RunCommand : コマンドを実行
func (a *App) RunCommand(cmd string) error {
	args, err := split(cmd)
	if err != nil {
		return err
	}

	return a.cmd.Execute(args)
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
			if err := a.RunCommand(cmd); err != nil {
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
		shared.RequestExecCommand("docs shortcuts")
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
