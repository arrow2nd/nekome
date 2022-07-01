package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
	a.cmd.CompletionOptions.HiddenDefaultCmd = !shared.isCommandLineMode

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
	a.cmd.SetHelpFunc(func(c *cobra.Command, s []string) {
		help := c.Long
		if help == "" {
			help = c.Short
		}

		text := fmt.Sprintf("%s\n\n%s", help, c.UsageString())
		a.view.AddPage(newHelpPage(c.Name(), text), true)
	})
}
