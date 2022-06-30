package app

import "github.com/spf13/cobra"

// newCmd : コマンド生成
func newCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "nekome",
		Short: "TUI Twitter client 🐈",
		Long:  "nekome is a TUI Twitter client that runs on the terminal 🐈",
	}
}

// initCmd : コマンド初期化
func (a *App) initCmd() {
	a.cmd.SilenceUsage = true
	a.cmd.SilenceErrors = true

	a.cmd.AddCommand(
		a.newHomeCmd(),
		a.newMentionCmd(),
		a.newListCmd(),
		a.newUserCmd(),
		a.newSearchCmd(),
		a.newTweetCmd(),
		a.newQuitCmd(),
	)
}
