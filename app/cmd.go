package app

import (
	"fmt"

	"github.com/arrow2nd/nekome/cli"
	"github.com/arrow2nd/nekome/log"
	"github.com/spf13/pflag"
)

// newCmd : コマンド生成
func newCmd() *cli.Command {
	return &cli.Command{
		Name:  "nekome",
		Short: "TUI Twitter client 🐈",
		Long:  "nekome is a TUI Twitter client that runs on the terminal 🐈",
		SetFlag: func(f *pflag.FlagSet) {
			f.StringP("user", "u", shared.conf.Settings.Feature.MainUser, "specify user to use")
			f.BoolP("version", "v", false, "show version")
		},
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			ver, _ := f.GetBool("version")

			// バージョンを表示
			if shared.isCommandLineMode && ver {
				log.LogExit(fmt.Sprintf("🐈 nekome for v.%s", version))
			}

			return nil
		},
	}
}

// initCmd : コマンド初期化
func (a *App) initCmd() {
	// コマンド追加
	a.cmd.AddCommand(
		a.newHomeCmd(),
		a.newMentionCmd(),
		a.newListCmd(),
		a.newUserCmd(),
		a.newSearchCmd(),
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
