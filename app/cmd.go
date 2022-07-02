package app

import "github.com/arrow2nd/nekome/cli"

// newCmd : コマンド生成
func newCmd() *cli.Command {
	return &cli.Command{
		Name:  "nekome",
		Short: "TUI Twitter client 🐈",
		Long:  "nekome is a TUI Twitter client that runs on the terminal 🐈",
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
		a.newHelpShortCuts(),
	)

	if shared.isCommandLineMode {
		return
	}

	// ヘルプの出力を新規ページに割り当てる
	a.cmd.Help = func(c *cli.Command, h string) {
		a.view.AddPage(newHelpPage(c.Name, h), true)
	}
}
