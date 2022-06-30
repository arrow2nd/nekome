package app

import "github.com/spf13/cobra"

// newQuitCmd : quitコマンド生成
func (a *App) newQuitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "quit",
		Aliases: []string{"q"},
		Short:   "Quit the application",
		Args:    cobra.NoArgs,
		Hidden:  shared.isCommandLineMode,
		Run: func(cmd *cobra.Command, args []string) {
			a.quitApp()
		},
	}

	return cmd
}
