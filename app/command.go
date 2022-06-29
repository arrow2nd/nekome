package app

import (
	"errors"
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"
)

// ExecCommand : コマンドを実行
func (a *App) ExecCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("command not found")
	}

	// フラグが無い・独自のフラグがあるコマンド
	switch args[0] {
	case "tweet", "t":
		return a.postTweet(args)
	case "quit", "q":
		a.quitApp()
		return nil
	}

	// フラグを設定
	var unfocus bool
	f := flag.NewFlagSet("nekome", flag.ContinueOnError)
	f.BoolVarP(&unfocus, "unfocus", "u", false, "")

	// フラグをパース
	if err := f.Parse(args); err != nil {
		return err
	}

	// ページ系のコマンド
	switch f.Arg(0) {
	case "home", "h":
		return a.view.AddPage(newTimelinePage(homeTL), !unfocus)
	case "mention", "m":
		return a.view.AddPage(newTimelinePage(mentionTL), !unfocus)
	case "list", "l":
		return a.openListPage(f.Arg(1), f.Arg(2), !unfocus)
	case "user", "u":
		return a.openUserPage(f.Arg(1), !unfocus)
	case "search", "s":
		return a.openSearchPage(f.Arg(1), !unfocus)
	}

	return fmt.Errorf(`"%s" is not a command`, f.Arg(0))
}

// openListPage : リストページを開く
func (a *App) openListPage(name, id string, focus bool) error {
	if name == "" {
		return errors.New("please specify list name")
	}

	if id == "" {
		return errors.New("please specify list id")
	}

	return a.view.AddPage(newListPage(name, id), focus)

}

// openUserPage : ユーザページを開く
func (a *App) openUserPage(userName string, focus bool) error {
	// ユーザの指定がないなら自分を指定
	if userName == "" {
		userName = shared.api.CurrentUser.UserName
	}

	// @を除去
	userName = strings.ReplaceAll(userName, "@", "")

	return a.view.AddPage(newUserPage(userName), focus)
}

// openSearchPage : 検索ページを開く
func (a *App) openSearchPage(query string, focus bool) error {
	if query == "" {
		return errors.New("please specify search keywords")
	}

	return a.view.AddPage(newSearchPage(query), focus)
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
