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
		"user",
		"list",
		"search",
		"switch",
		"exit",
	}
}

// ExecCmd : コマンドを実行
func (a *App) ExecCmd(args []string) error {
	// tweetコマンドのフラグを設定
	tweetFlag := flag.NewFlagSet("tweet", flag.ContinueOnError)

	tweetFlag.AddFlag(&flag.Flag{
		Name:      "quote",
		Shorthand: "q",
		Usage:     "Specify the ID of the tweet to quote",
	})

	tweetFlag.AddFlag(&flag.Flag{
		Name:      "reply",
		Shorthand: "r",
		Usage:     "Specify the ID of the tweet to which you are replying",
	})

	// 引数をパース
	f := flag.NewFlagSet("nekome", flag.ContinueOnError)
	f.Parse(args)

	if f.NArg() == 0 {
		return errors.New("command not found")
	}

	switch f.Arg(0) {
	case "user", "u":
		return a.cmdOpenUserPage(f.Arg(1))
	}

	return fmt.Errorf(`"%s" is not a command`, f.Arg(0))
}

// cmdOpenUserPage : ユーザページを開く
func (a *App) cmdOpenUserPage(userName string) error {
	// ユーザの指定がないなら自分を指定
	if userName == "" {
		userName = shared.api.CurrentUser.UserName
	}

	// @を除去
	userName = strings.ReplaceAll(userName, "@", "")

	return a.view.AddPage(newUserPage(userName), true)
}
