package app

import (
	"github.com/arrow2nd/nekome/cli"
	"github.com/spf13/pflag"
)

// newQuitCmd : quitコマンド生成
func (a *App) newQuitCmd() *cli.Command {
	return &cli.Command{
		Name:      "quit",
		Shorthand: "q",
		Short:     "Quit the application",
		Validate:  cli.NoArgs(),
		Hidden:    shared.isCommandLineMode,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			a.quitApp()
			return nil
		},
	}
}
