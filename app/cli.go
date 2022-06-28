package app

import (
	"errors"
	"fmt"
	"os"
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

// postTweet : ツイートを投稿
func (a *App) postTweet(args []string) error {
	quote := ""
	reply := ""
	editor := ""

	// フラグを設定
	f := flag.NewFlagSet("tweet", flag.ContinueOnError)
	f.StringVarP(&quote, "quote", "q", "", "Specify the ID of the tweet to quote")
	f.StringVarP(&reply, "reply", "r", "", "Specify the ID of the tweet to which you are replying")
	f.StringVarP(&editor, "editor", "e", os.Getenv("EDITOR"), "Specify the editor to start (Default is $EDITOR)")

	if err := f.Parse(args); err != nil {
		return err
	}

	text := f.Arg(1)

	// エディタを起動して編集
	if text == "" {
		t, err := editTextInEditor(editor)
		if err != nil {
			return err
		}

		text = t

		// 画面が鬼のように崩れるので再描画
		a.app.Sync()
	}

	// ツイート末尾の改行を削除
	text = strings.TrimRight(text, "\n")
	if strings.HasSuffix(text, "\r") {
		text = strings.TrimRight(text, "\r")
	}

	// 投稿
	post := func() {
		if err := shared.api.PostTweet(text, quote, reply); err != nil {
			shared.SetErrorStatus("Tweet", err.Error())
			return
		}

		shared.SetStatus("Tweeted", text)
	}

	// 確認画面が不要
	if !shared.conf.Settings.Feature.Confirm["Tweet"] {
		post()
		return nil
	}

	shared.ReqestPopupModal(&ModalOpt{
		title:  fmt.Sprintf("Do you want to tweet?\n\n\"%s\"", text),
		onDone: post,
	})

	return nil
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
