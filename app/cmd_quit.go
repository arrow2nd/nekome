package app

import (
	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
)

// newQuitCmd : quitコマンド生成
func (a *App) newQuitCmd() *cli.Command {
	return &cli.Command{
		Name:         "quit",
		Alias:        "q",
		Short:        "Quit the application",
		ValidateFunc: cli.NoArgs(),
		Hidden:       shared.isCommandLineMode,
		RunFunc: func(c *cli.Command, f *pflag.FlagSet) error {
			a.quitApp()
			return nil
		},
	}
}
