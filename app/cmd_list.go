package app

import (
	"errors"

	"github.com/arrow2nd/nekome/v2/cli"
	"github.com/spf13/pflag"
)

func (a *App) newListCmd() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Shorthand: "l",
		Short:     "Add list timeline page",
		UsageArgs: "<list name> <list ID>",
		Example:   "list cathouse 1234567890",
		Validate:  cli.RequireArgs(2),
		Hidden:    shared.isCommandLineMode,
		SetFlag:   setUnfocusFlag,
		Run: func(c *cli.Command, f *pflag.FlagSet) error {
			name := f.Arg(0)
			if name == "" {
				return errors.New("please specify list name")
			}

			id := f.Arg(1)
			if id == "" {
				return errors.New("please specify list id")
			}

			page, err := newListPage(name, id)
			if err != nil {
				return err
			}

			unfocus, _ := f.GetBool("unfocus")
			if err := a.view.AddPage(page, !unfocus); err != nil {
				return err
			}

			go page.Load()

			return nil
		},
	}
}
