package app

import (
	"errors"
	"os"

	"github.com/arrow2nd/nekome/cli"
	"github.com/arrow2nd/nekome/config"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

var version = "develop"

// App : アプリケーション
type App struct {
	app                   *tview.Application
	cmd                   *cli.Command
	view                  *view
	statusBar             *statusBar
	commandLine           *commandLine
	isDisablePageKeyEvent bool
}

// New : 新規作成
func New() *App {
	return &App{
		app:                   tview.NewApplication(),
		cmd:                   newCmd(),
		view:                  nil,
		statusBar:             nil,
		commandLine:           nil,
		isDisablePageKeyEvent: false,
	}
}

// Init : 初期化
func (a *App) Init() error {
	// 設定読み込み
	if err := a.loadConfig(); err != nil {
		return err
	}

	// 実行時の引数をパース
	user, err := a.parseRuntimeArgs()
	if err != nil {
		return err
	}

	// ログイン
	if user != "" {
		if err := loginAccount(user); err != nil {
			return err
		}
	} else {
		return errors.New(
			"main user is not set, please run 'nekome edit' and set it to 'feature.main_user' in settings.toml",
		)
	}

	// コマンド初期化
	a.initCmd()

	// コマンドラインモードならUIの初期化をスキップ
	if shared.isCommandLineMode {
		return nil
	}

	// UI準備
	a.setAppStyles()
	a.view = newView()
	a.statusBar = newStatusBar()
	a.commandLine = newCommandLine()

	// 日本語環境等での罫線の乱れ対策
	// LINK: https://github.com/mattn/go-runewidth/issues/14
	runewidth.DefaultCondition.EastAsianWidth = !shared.conf.Settings.Feature.IsLocaleCJK

	// ページのキーハンドラを設定
	a.view.SetInputCapture(a.handlePageKeyEvent)

	// ステータスバー初期化
	a.statusBar.Init()
	a.statusBar.DrawAccountInfo()

	// コマンドライン初期化
	a.commandLine.Init()
	a.initAutocomplate()

	// 画面レイアウト
	layout := tview.NewGrid().
		SetRows(1, 0, 1, 1).
		SetBorders(false).
		AddItem(a.view.tabView, 0, 0, 1, 1, 0, 0, false).
		AddItem(a.statusBar.flex, 2, 0, 1, 1, 0, 0, false).
		AddItem(a.commandLine.inputField, 3, 0, 1, 1, 0, 0, false).
		AddItem(a.view.mainView, 1, 0, 1, 1, 0, 0, true)

	a.app.
		SetRoot(layout, true).
		SetInputCapture(a.handleGlobalKeyEvents)

	a.runStartupCommands()

	return nil
}

// setAppStyles : アプリ全体のスタイルを設定
func (a *App) setAppStyles() {
	app := shared.conf.Style.App

	bgColor := app.BackgroundColor.ToColor()
	textColor := app.TextColor.ToColor()
	borderColor := app.BorderColor.ToColor()

	// 背景色
	tview.Styles.PrimitiveBackgroundColor = bgColor
	tview.Styles.ContrastBackgroundColor = bgColor
	tview.Styles.MoreContrastBackgroundColor = app.BackgroundColor.ToColor()

	// テキスト色
	tview.Styles.PrimaryTextColor = textColor
	tview.Styles.ContrastSecondaryTextColor = textColor
	tview.Styles.TitleColor = textColor
	tview.Styles.TertiaryTextColor = app.SubTextColor.ToColor()

	// ボーダー色
	tview.Styles.BorderColor = borderColor
	tview.Styles.GraphicsColor = borderColor
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
	existUser, err := shared.conf.LoadCred()
	if err != nil {
		return err
	}

	if existUser {
		return nil
	}

	// ユーザ情報が無い場合、新規追加
	return addAccount(true)
}

// parseRuntimeArgs : 実行時の引数をパースして、ログインユーザを返す
func (a *App) parseRuntimeArgs() (string, error) {
	f := a.cmd.NewFlagSet()

	if err := f.Parse(os.Args[1:]); err != nil {
		return "", err
	}

	// ログイン処理をスキップするか
	skipLogin := f.Changed("help") || f.Changed("version") || f.Arg(0) == "e" || f.Arg(0) == "edit"

	// コマンドラインモードか
	shared.isCommandLineMode = f.NArg() > 0 || skipLogin

	// ログインしない
	if skipLogin {
		return "", nil
	}

	user, _ := f.GetString("user")
	return user, nil
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
	for _, c := range shared.conf.Settings.Feature.StartupCmds {
		if err := a.ExecCommnad(c); err != nil {
			shared.SetErrorStatus("Command", err.Error())
		}
	}
}

// Run : 実行
func (a *App) Run() error {
	// コマンドラインモード
	if shared.isCommandLineMode {
		return a.cmd.Execute(os.Args[1:])
	}

	go a.eventReciever()
	return a.app.Run()
}

// ExecCommnad : コマンドを実行
func (a *App) ExecCommnad(cmd string) error {
	args, err := split(cmd)
	if err != nil {
		return err
	}

	return a.cmd.Execute(args)
}

// quitApp : アプリを終了
func (a *App) quitApp() {
	// 確認画面が不要ならそのまま終了
	if !shared.conf.Settings.Confirm["quit"] {
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
		case b := <-shared.chDisablePageKeyEvent:
			a.isDisablePageKeyEvent = b
		case opt := <-shared.chPopupModal:
			a.view.PopupModal(opt)
			a.app.Draw()
		case cmd := <-shared.chExecCommand:
			if err := a.ExecCommnad(cmd); err != nil {
				shared.SetErrorStatus("Command", err.Error())
			}
		case cmd := <-shared.chInputCommand:
			a.app.SetFocus(a.commandLine.inputField)
			a.commandLine.SetText(cmd)
			a.app.Draw()
		case <-shared.chFocusMainView:
			focus := a.app.GetFocus()
			if focus != a.view.textArea {
				a.app.SetFocus(a.view.mainView)
			}
			a.app.Draw()
		case p := <-shared.chFocusPrimitive:
			a.app.SetFocus(*p)
			a.app.Draw()
		}
	}
}

// handleGlobalKeyEvents : アプリ全体のキーハンドラ
func (a *App) handleGlobalKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	if a.isDisablePageKeyEvent {
		return event
	}

	// アプリを終了
	if event.Key() == tcell.KeyCtrlQ {
		a.quitApp()
		return nil
	}

	return event
}

// handlePageKeyEvent : ページビューのキーハンドラ
func (a *App) handlePageKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	if a.isDisablePageKeyEvent {
		return event
	}

	key := event.Key()
	keyRune := event.Rune()

	// 前のタブを選択
	if key == tcell.KeyLeft || keyRune == 'h' {
		a.view.MoveTab(TabMovePrev)
		return nil
	}

	// 次のタブを選択
	if key == tcell.KeyRight || keyRune == 'l' {
		a.view.MoveTab(TabMoveNext)
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

	// キーバインドのヘルプ
	if keyRune == '?' {
		shared.RequestExecCommand("docs keybindings")
		return nil
	}

	// ページを削除
	if key == tcell.KeyCtrlW {
		a.view.RemoveCurrentPage()
		return nil
	}

	return event
}
