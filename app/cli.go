package app

import (
	"errors"
	"fmt"
	"strings"

	flag "github.com/spf13/pflag"
)

// getCommands : コマンドリストを取得
func (a *App) getCommands() []string {
	return []string{
		"tweet",
		"home",
		"mention",
		"user",
		"list",
		"search",
		"switch",
		"quit",
	}
}

// ExecCmd : コマンドを実行
func (a *App) ExecCmd(args []string) error {
	// フラグを設定
	f := flag.NewFlagSet("nekome", flag.ContinueOnError)
	tweetFlag := flag.NewFlagSet("tweet", flag.ContinueOnError)

	f.BoolP("unfocus", "u", false, "")
	tweetFlag.BoolP("quote", "q", false, "Specify the ID of the tweet to quote")
	tweetFlag.BoolP("reply", "r", false, "Specify the ID of the tweet to which you are replying")

	// 引数をパース
	if err := f.Parse(args); err != nil {
		return err
	}

	if f.NArg() == 0 {
		return errors.New("command not found")
	}

	unfocus, _ := f.GetBool("unfocus")

	// コマンドを解析
	switch f.Arg(0) {
	case "home", "h":
		return a.view.AddPage(newTimelinePage(homeTL), !unfocus)
	case "mention", "m":
		return a.view.AddPage(newTimelinePage(mentionTL), !unfocus)
	case "user", "u":
		return a.openUserPage(f.Arg(1), !unfocus)
	case "quit", "q":
		a.quitApp()
		return nil
	}

	return fmt.Errorf(`"%s" is not a command`, f.Arg(0))
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

// quitApp : アプリを終了
func (a *App) quitApp() {
	a.blurCommandLine()

	a.view.PopupModal(&ModalOpt{
		title:  "Do you want to quit the app?",
		onDone: a.app.Stop,
	})

}
