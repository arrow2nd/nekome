package app

import "github.com/spf13/cobra"

func (a *App) newQuitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "quit",
		Long: "Quit the application",
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			a.quitApp()
		},
	}

	return cmd
}
