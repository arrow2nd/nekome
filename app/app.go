package app

import (
	"os"

	"code.rocketnine.space/tslocum/cbind"
	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/arrow2nd/nekome/v2/config"
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
	if err := a.loadConfig(); err != nil {
		return err
	}

	isLoginSkip, user, err := a.parseRuntimeArgs()
	if err != nil {
		return err
	}

	if !isLoginSkip {
		if err := loginAccount(user); err != nil {
			return err
		}
	}

	a.initCommands()

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
	runewidth.DefaultCondition.EastAsianWidth = !shared.conf.Pref.Feature.IsLocaleCJK

	// Ctrl+K/Jの再マッピングを無効化
	cbind.UnifyEnterKeys = false

	// キーバインドを設定
	if err := a.setGlobalKeybindings(); err != nil {
		return err
	}
	if err := a.setViewKeybindings(); err != nil {
		return err
	}

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
		AddItem(a.view.tabBar, 0, 0, 1, 1, 0, 0, false).
		AddItem(a.statusBar.flex, 2, 0, 1, 1, 0, 0, false).
		AddItem(a.commandLine.inputField, 3, 0, 1, 1, 0, 0, false).
		AddItem(a.view.flex, 1, 0, 1, 1, 0, 0, true)

	a.app.SetRoot(layout, true)

	a.execStartupCommands()

	return nil
}

// loadConfig : 設定を読み込む
func (a *App) loadConfig() error {
	shared.conf = config.New()

	shared.conf.CheckOldFile()

	// 環境設定
	if err := shared.conf.LoadPreferences(); err != nil {
		return err
	}

	// スタイル定義
	if err := shared.conf.LoadStyle(); err != nil {
		return err
	}

	// 認証情報
	return shared.conf.LoadCred()
}

// parseRuntimeArgs : 実行時の引数をパースして、ログインユーザを返す
func (a *App) parseRuntimeArgs() (bool, string, error) {
	f := a.cmd.NewFlagSet()

	f.ParseErrorsWhitelist.UnknownFlags = true

	if err := f.Parse(os.Args[1:]); err != nil {
		return false, "", err
	}

	// ログインをスキップするか
	arg := f.Arg(0)
	isSkipLogin := f.Changed("help") || f.Changed("version") || arg == "e" || arg == "edit"

	// コマンドラインモードか
	shared.isCommandLineMode = f.NArg() > 0 || isSkipLogin

	if isSkipLogin {
		return true, "", nil
	}

	user, _ := f.GetString("user")
	return false, user, nil
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

	// ボーダー
	tview.Borders.HorizontalFocus = tview.BoxDrawingsHeavyHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsHeavyVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsHeavyDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsHeavyDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsHeavyUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsHeavyUpAndLeft
}

// setGlobalKeybindings : アプリ全体のキーバインドを設定
func (a *App) setGlobalKeybindings() error {
	handlers := map[string]func(){
		config.ActionQuit: func() {
			a.quitApp()
		},
	}

	c, err := shared.conf.Pref.Keybindings.Global.MappingEventHandler(handlers)
	if err != nil {
		return err
	}

	a.app.SetInputCapture(a.warpKeyEventHandler(c))

	return nil
}

// setViewKeybindings : ビューのキーバインドを設定
func (a *App) setViewKeybindings() error {
	handlers := map[string]func(){
		config.ActionSelectPrevTab: func() {
			a.view.MoveTab(TabMovePrev)
		},
		config.ActionSelectNextTab: func() {
			a.view.MoveTab(TabMoveNext)
		},
		config.ActionRedraw: func() {
			a.app.Sync()
		},
		config.ActionFocusCmdLine: func() {
			a.app.SetFocus(a.commandLine.inputField)
		},
		config.ActionShowHelp: func() {
			shared.RequestExecCommand("docs keybindings")
		},
		config.ActionClosePage: func() {
			a.view.CloseCurrentPage()
		},
	}

	c, err := shared.conf.Pref.Keybindings.View.MappingEventHandler(handlers)
	if err != nil {
		return err
	}

	a.view.SetInputCapture(a.warpKeyEventHandler(c))

	return nil
}

// warpKeyEventHandler : イベントハンドラのラップ関数
func (a *App) warpKeyEventHandler(c *cbind.Configuration) func(*tcell.EventKey) *tcell.EventKey {
	return func(ev *tcell.EventKey) *tcell.EventKey {
		// 操作が無効
		if a.isDisablePageKeyEvent {
			return ev
		}

		return c.Capture(ev)
	}
}

// initCommands : コマンドを初期化
func (a *App) initCommands() {
	a.cmd.AddCommand(
		a.newHomeCmd(),
		a.newMentionCmd(),
		a.newListCmd(),
		a.newUserCmd(),
		a.newSearchCmd(),
		a.newLikesCmd(),
		a.newTweetCmd(),
		a.newQuitCmd(),
		a.newDocsCmd(),
		a.newAccountCmd(),
		a.newEditCmd(),
	)

	if shared.isCommandLineMode {
		return
	}

	// ヘルプの出力を新規ページに割り当てる
	a.cmd.Help = func(c *cli.Command, h string) {
		a.view.AddPage(newDocsPage(c.Name, h), true)
	}
}

// initAutocomplate : 入力補完を初期化
func (a *App) initAutocomplate() {
	cmds := a.cmd.GetChildrenNames(true)

	if err := a.commandLine.SetAutocompleteItems(cmds); err != nil {
		shared.SetErrorStatus("Init - CommandLine", err.Error())
	}

}

// execStartupCommands : 起動時に実行するコマンドを一括で実行
func (a *App) execStartupCommands() {
	for _, c := range shared.conf.Pref.Feature.StartupCmds {
		if err := a.ExecCommnad(c); err != nil {
			shared.SetErrorStatus("Command", err.Error())
		}
	}
}

// ExecCommnad : コマンドを実行
func (a *App) ExecCommnad(cmd string) error {
	args, err := split(cmd)
	if err != nil {
		return err
	}

	return a.cmd.Execute(args)
}

// Run : アプリを実行
func (a *App) Run() error {
	// コマンドラインモード
	if shared.isCommandLineMode {
		return a.cmd.Execute(os.Args[1:])
	}

	go a.eventReciever()

	return a.app.Run()
}

// eventReciever : イベントレシーバ
func (a *App) eventReciever() {
	for {
		select {
		case status := <-shared.chStatus:
			// ステータスメッセージを表示
			a.commandLine.UpdateStatusMessage(status)
			a.app.Draw()
		case indicator := <-shared.chIndicator:
			// インジケータを更新
			a.statusBar.DrawPageIndicator(indicator)
			a.app.Draw()
		case b := <-shared.chDisableViewKeyEvent:
			// ビューのキー操作ロック状態を更新
			a.isDisablePageKeyEvent = b
		case opt := <-shared.chPopupModal:
			// モーダルを表示
			a.view.PopupModal(opt)
			a.app.Draw()
		case cmd := <-shared.chExecCommand:
			// コマンドを実行`
			if err := a.ExecCommnad(cmd); err != nil {
				shared.SetErrorStatus("Command", err.Error())
			}
		case cmd := <-shared.chInputCommand:
			// コマンドを入力
			a.app.SetFocus(a.commandLine.inputField)
			a.commandLine.SetText(cmd)
			a.app.Draw()
		case <-shared.chFocusView:
			// ビューにフォーカス
			if a.app.GetFocus() != a.view.textArea {
				a.app.SetFocus(a.view.flex)
			}
			a.app.Draw()
		case p := <-shared.chFocusPrimitive:
			// 任意のプリミティブにフォーカス
			a.app.SetFocus(*p)
			a.app.Draw()
		}
	}
}

// quitApp : アプリを終了
func (a *App) quitApp() {
	// 確認画面が不要ならそのまま終了
	if !shared.conf.Pref.Confirm["quit"] {
		a.app.Stop()
		return
	}

	a.view.PopupModal(&ModalOpt{
		title:  "Do you want to quit the app?",
		onDone: a.app.Stop,
	})
}
